package io

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDbParser struct{}

func (p *MariaDbParser) TargetName() string {
	return "MariaDB"
}

// ------------------------
// DISCOVERY IMPLEMENTATION
// ------------------------

func (p *MariaDbParser) GetDatabaseList(host string, port int, user, password string) ([]string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		databases = append(databases, name)
	}

	return databases, nil
}

func (p *MariaDbParser) GetTableList(host string, port int, user, password, database string) ([]string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}

	return tables, nil
}

func (p *MariaDbParser) GetTableSchema(host string, port int, user, password, database, tableName string) (*Entity, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			COLUMN_NAME, 
			DATA_TYPE, 
			COLUMN_TYPE, 
			IS_NULLABLE, 
			COLUMN_DEFAULT, 
			CHARACTER_MAXIMUM_LENGTH, 
			NUMERIC_PRECISION, 
			COLUMN_KEY 
		FROM information_schema.columns 
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ORDINAL_POSITION`

	rows, err := db.Query(query, database, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []Field
	var primaryKeys []string

	for rows.Next() {
		var colName, dataType, colType, isNullable, colKey string
		var colDefault, charLen, numPrec sql.NullString

		err := rows.Scan(&colName, &dataType, &colType, &isNullable, &colDefault, &charLen, &numPrec, &colKey)
		if err != nil {
			return nil, err
		}

		colMap := map[string]any{
			"COLUMN_NAME":              colName,
			"DATA_TYPE":                dataType,
			"COLUMN_TYPE":              colType,
			"IS_NULLABLE":              isNullable,
			"COLUMN_DEFAULT":           nil,
			"CHARACTER_MAXIMUM_LENGTH": nil,
			"COLUMN_KEY":               colKey,
		}
		if colDefault.Valid {
			colMap["COLUMN_DEFAULT"] = colDefault.String
		}
		if charLen.Valid {
			colMap["CHARACTER_MAXIMUM_LENGTH"] = charLen.String
		} else if numPrec.Valid {
			colMap["CHARACTER_MAXIMUM_LENGTH"] = numPrec.String
		}

		field := p.mapColumn(colMap)
		fields = append(fields, field)

		if colKey == "PRI" {
			primaryKeys = append(primaryKeys, colName)
		}
	}

	return &Entity{
		Name:       tableName,
		Fields:     fields,
		PrimaryKey: primaryKeys,
	}, nil
}

// ----------------------------------------
// EXPORT: Intermediate Schema -> MariaDB DDL
// ----------------------------------------

func (p *MariaDbParser) ToTarget(schema IntermediateSchema) string {
	var sb strings.Builder

	for _, entity := range schema.Schema.Entities {
		sb.WriteString(p.generateTable(entity))
		sb.WriteString("\n\n")
	}

	for _, rel := range schema.Schema.Relationships {
		sb.WriteString(p.generateForeignKey(rel))
		sb.WriteString("\n\n")
	}

	return strings.TrimSpace(sb.String())
}

func (p *MariaDbParser) generateTable(entity Entity) string {
	var columns []string
	for _, field := range entity.Fields {
		columns = append(columns, p.generateColumn(field))
	}

	primaryKeyClause := ""
	if len(entity.PrimaryKey) > 0 {
		var pkCols []string
		for _, k := range entity.PrimaryKey {
			pkCols = append(pkCols, fmt.Sprintf("`%s`", k))
		}
		primaryKeyClause = fmt.Sprintf(",\n  PRIMARY KEY (%s)", strings.Join(pkCols, ", "))
	}

	createTable := fmt.Sprintf("CREATE TABLE `%s` (\n%s%s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		entity.Name,
		strings.Join(columns, ",\n"),
		primaryKeyClause,
	)

	var indexStatements []string
	if entity.Indexes != nil {
		for _, idx := range entity.Indexes {
			unique := ""
			if idx.Unique != nil && *idx.Unique {
				unique = "UNIQUE "
			}
			var fields []string
			for _, f := range idx.Fields {
				fields = append(fields, fmt.Sprintf("`%s`", f))
			}
			indexStatements = append(indexStatements,
				fmt.Sprintf("CREATE %sINDEX `%s` ON `%s` (%s);",
					unique, idx.Name, entity.Name, strings.Join(fields, ", ")),
			)
		}
	}

	if len(indexStatements) > 0 {
		return createTable + "\n" + strings.Join(indexStatements, "\n")
	}
	return createTable
}

func (p *MariaDbParser) generateColumn(field Field) string {
	typeStr := p.mapFieldType(field)
	nullable := ""
	if !field.Nulish {
		nullable = " NOT NULL"
	}
	defaultClause := p.getDefaultClause(field)

	return fmt.Sprintf("  `%s` %s%s%s", field.Name, typeStr, nullable, defaultClause)
}

func (p *MariaDbParser) mapFieldType(field Field) string {
	if len(field.EnumValues) > 0 {
		var quoted []string
		for _, v := range field.EnumValues {
			quoted = append(quoted, fmt.Sprintf("'%s'", v))
		}
		return fmt.Sprintf("ENUM(%s)", strings.Join(quoted, ", "))
	}

	switch field.Type {
	case TypeString:
		return "TEXT"
	case TypeInteger:
		return "INT"
	case TypeFloat:
		return "DOUBLE"
	case TypeBoolean:
		return "TINYINT(1)"
	case TypeDateTime:
		return "DATETIME(6)"
	case TypeDate:
		return "DATE"
	case TypeObject, TypeArray, TypeMap:
		return "JSON"
	default:
		return "TEXT"
	}
}

func (p *MariaDbParser) getDefaultClause(field Field) string {
	if field.Default == nil {
		return ""
	}

	if s, ok := field.Default.(string); ok && s == "now()" {
		if field.Type == TypeDateTime {
			return " DEFAULT CURRENT_TIMESTAMP(6)"
		}
		return " DEFAULT CURRENT_TIMESTAMP"
	}

	if field.Type == TypeBoolean {
		val := "0"
		if b, ok := field.Default.(bool); ok && b {
			val = "1"
		}
		return " DEFAULT " + val
	}

	switch v := field.Default.(type) {
	case string:
		return fmt.Sprintf(" DEFAULT '%s'", v)
	case int, int64, float64:
		return fmt.Sprintf(" DEFAULT %v", v)
	}

	if field.Type == TypeObject || field.Type == TypeArray || field.Type == TypeMap {
		return " DEFAULT ('{}')"
	}

	return ""
}

func (p *MariaDbParser) generateForeignKey(rel Relationship) string {
	onDelete := p.mapReferentialAction(rel.OnDelete)
	onUpdate := p.mapReferentialAction(rel.OnUpdate)

	return fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `fk_%s_%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`) ON DELETE %s ON UPDATE %s;",
		rel.Source.Entity,
		rel.Source.Entity,
		rel.Source.Field,
		rel.Source.Field,
		rel.Target.Entity,
		rel.Target.Field,
		onDelete,
		onUpdate,
	)
}

func (p *MariaDbParser) mapReferentialAction(action ReferentialAction) string {
	switch action {
	case ActionCascade:
		return "CASCADE"
	case ActionRestrict:
		return "RESTRICT"
	case ActionSetNull:
		return "SET NULL"
	case ActionNoAction:
		return "NO ACTION"
	default:
		return "RESTRICT"
	}
}

// ----------------------------------------
// IMPORT: MariaDB -> Intermediate Schema
// ----------------------------------------

type IntrospectionData struct {
	Tables      []map[string]any `json:"tables"`
	Columns     []map[string]any `json:"columns"`
	Indexes     []map[string]any `json:"indexes"`
	ForeignKeys []map[string]any `json:"foreignKeys"`
}

func (p *MariaDbParser) FromTarget(source IntrospectionData) IntermediateSchema {
	entities := []Entity{}
	relationships := []Relationship{}

	columnsByTable := make(map[string][]map[string]any)
	for _, col := range source.Columns {
		tableName := col["TABLE_NAME"].(string)
		columnsByTable[tableName] = append(columnsByTable[tableName], col)
	}

	for _, tableMap := range source.Tables {
		tableName := tableMap["TABLE_NAME"].(string)
		cols := columnsByTable[tableName]

		entityIndexes := []Index{}
		for _, i := range source.Indexes {
			if i["TABLE_NAME"] == tableName {
				unique := i["NON_UNIQUE"].(int64) == 0
				idx := Index{
					Name:   i["INDEX_NAME"].(string),
					Fields: i["COLUMNS"].([]string),
					Unique: &unique,
				}
				entityIndexes = append(entityIndexes, idx)
			}
		}

		fields := []Field{}
		for _, col := range cols {
			fields = append(fields, p.mapColumn(col))
		}

		entity := Entity{
			Name:        tableName,
			Description: "",
			PrimaryKey:  p.getPrimaryKey(cols),
			Fields:      fields,
		}
		if len(entityIndexes) > 0 {
			entity.Indexes = entityIndexes
		}

		entities = append(entities, entity)
	}

	for _, fk := range source.ForeignKeys {
		rel := Relationship{
			Name: fmt.Sprintf("fk_%s", fk["CONSTRAINT_NAME"]),
			Type: "belongsTo",
			Source: SourceTarget{
				Entity: fk["TABLE_NAME"].(string),
				Field:  fk["COLUMN_NAME"].(string),
			},
			Target: SourceTarget{
				Entity: fk["REFERENCED_TABLE_NAME"].(string),
				Field:  fk["REFERENCED_COLUMN_NAME"].(string),
			},
			OnDelete: p.mapFkAction(fk["DELETE_RULE"].(string)),
			OnUpdate: p.mapFkAction(fk["UPDATE_RULE"].(string)),
		}
		relationships = append(relationships, rel)
	}

	return IntermediateSchema{
		Meta: Meta{
			Version: "1.0",
			Database: Database{
				From: "MariaDB",
				Name: "imported_db",
			},
			Description: "Imported from MariaDB",
		},
		Schema: Schema{
			Entities:      entities,
			Relationships: relationships,
		},
	}
}

func (p *MariaDbParser) getPrimaryKey(columns []map[string]any) []string {
	var pk []string
	for _, col := range columns {
		if key, ok := col["COLUMN_KEY"].(string); ok && key == "PRI" {
			pk = append(pk, col["COLUMN_NAME"].(string))
		}
	}
	return pk
}

func (p *MariaDbParser) mapColumn(col map[string]any) Field {
	name := col["COLUMN_NAME"].(string)
	columnType := col["COLUMN_TYPE"].(string)
	isNullable := col["IS_NULLABLE"].(string) == "YES"
	columnDefault := col["COLUMN_DEFAULT"]

	var enumValues []string
	if strings.HasPrefix(columnType, "enum(") {
		enumValues = p.parseEnumFromColumnType(columnType)
	}

	var fieldType FieldType
	if len(enumValues) > 0 {
		fieldType = TypeString
	} else {
		typeName := col["DATA_TYPE"].(string)
		switch strings.ToLower(typeName) {
		case "char", "varchar", "text", "longtext":
			fieldType = TypeString
		case "int", "integer":
			fieldType = TypeInteger
		case "double", "float":
			fieldType = TypeFloat
		case "tinyint":
			if columnType == "tinyint(1)" {
				fieldType = TypeBoolean
			} else {
				fieldType = TypeInteger
			}
		case "datetime", "timestamp":
			fieldType = TypeDateTime
		case "date":
			fieldType = TypeDate
		case "json":
			fieldType = TypeObject
		default:
			fieldType = TypeString
		}
	}

	var defaultVal any
	if columnDefault != nil {
		s := fmt.Sprintf("%v", columnDefault)
		if s == "CURRENT_TIMESTAMP" || strings.HasPrefix(s, "CURRENT_TIMESTAMP") {
			defaultVal = "now()"
		} else if fieldType == TypeBoolean {
			defaultVal = (s == "1" || s == "b'1'")
		} else {
			defaultVal = columnDefault
		}
	}

	var length *int
	if l, ok := col["CHARACTER_MAXIMUM_LENGTH"]; ok && l != nil {
		switch v := l.(type) {
		case string:
			var li int
			fmt.Sscanf(v, "%d", &li)
			if li > 0 {
				length = &li
			}
		case int:
			length = &v
		case int64:
			li := int(v)
			length = &li
		}
	}

	embed := fieldType == TypeObject
	return Field{
		Name:       name,
		Type:       fieldType,
		Length:     length,
		Nulish:     isNullable,
		Default:    defaultVal,
		Embed:      &embed,
		EnumValues: enumValues,
	}
}

func (p *MariaDbParser) parseEnumFromColumnType(columnType string) []string {
	start := strings.Index(columnType, "(")
	end := strings.LastIndex(columnType, ")")
	if start == -1 || end == -1 {
		return []string{}
	}
	inner := columnType[start+1 : end]
	parts := strings.Split(inner, ",")
	var result []string
	re := regexp.MustCompile(`^'+|'+$`)
	for _, p := range parts {
		s := strings.TrimSpace(p)
		s = re.ReplaceAllString(s, "")
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func (p *MariaDbParser) mapFkAction(action string) ReferentialAction {
	switch strings.ToLower(action) {
	case "cascade":
		return ActionCascade
	case "set null":
		return ActionSetNull
	case "no action":
		return ActionNoAction
	default:
		return ActionRestrict
	}
}
