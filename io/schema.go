package io

// OnConflictStrategy defines how to handle conflicts during data insertion
type OnConflictStrategy string

const (
	ConflictFail   OnConflictStrategy = "fail"
	ConflictIgnore OnConflictStrategy = "ignore"
	ConflictUpdate OnConflictStrategy = "update"
)

// ReferentialAction defines the action to take on a referenced record when it is modified or deleted
type ReferentialAction string

const (
	ActionCascade  ReferentialAction = "cascade"
	ActionRestrict ReferentialAction = "restrict"
	ActionSetNull  ReferentialAction = "setNull"
	ActionNoAction ReferentialAction = "noAction"
)

// FieldType represents the data types supported in the intermediate schema
type FieldType string

const (
	TypeString   FieldType = "string"
	TypeInteger  FieldType = "integer"
	TypeFloat    FieldType = "float"
	TypeBoolean  FieldType = "boolean"
	TypeDateTime FieldType = "datetime"
	TypeDate     FieldType = "date"
	TypeObject   FieldType = "object"
	TypeArray    FieldType = "array"
	TypeMap      FieldType = "map"
)

// IntermediateSchema is the root structure for a database project schema
type IntermediateSchema struct {
	Meta   Meta   `json:"meta"`
	Schema Schema `json:"schema"`
}

// Meta contains project-level metadata
type Meta struct {
	Version     string   `json:"version"`
	Database    Database `json:"database"`
	Description string   `json:"description"`
}

// Database identifies the source and name of the database
type Database struct {
	From string `json:"from,omitempty"`
	Name string `json:"name"`
}

// Schema contains the structured entities and their relationships
type Schema struct {
	Entities          []Entity            `json:"entities"`
	ForeignKeys       []ForeignKeys       `json:"foreignKeys"`
	ReferencedColumns []ReferencedColumns `json:"referencedColumns"`
}

// Entity represents a database table/collection
type Entity struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	PrimaryKey      []string            `json:"primaryKey"`
	OnConflict      *OnConflictStrategy `json:"onConflict,omitempty"`
	Fields          []Field             `json:"fields"`
	Indexes         []Index             `json:"indexes,omitempty"`
	ForeignEntities []ForeignKeys       `json:"ForeignKeys"`            // Foreign keys in this table
	ReferencedIn    []Reference         `json:"referencedIn,omitempty"` // Tables that reference this table
}

// ForeignEntity represents a foreign key relationship where this table references another
type ForeignEntity struct {
	ConstraintName  string            `json:"constraintName"`
	LocalField      string            `json:"localField"`      // Field in this table
	ReferencedTable string            `json:"referencedTable"` // Table being referenced
	ReferencedField string            `json:"referencedField"` // Field in referenced table
	OnDelete        ReferentialAction `json:"onDelete"`
	OnUpdate        ReferentialAction `json:"onUpdate"`
}

// Reference represents where this table is referenced by another table
type Reference struct {
	ConstraintName string            `json:"constraintName"`
	SourceTable    string            `json:"sourceTable"` // Table that references this one
	SourceField    string            `json:"sourceField"` // Field in source table
	LocalField     string            `json:"localField"`  // Field in this table being referenced
	OnDelete       ReferentialAction `json:"onDelete"`
	OnUpdate       ReferentialAction `json:"onUpdate"`
}

// Field represents a column/field in an entity
type Field struct {
	Name         string    `json:"name"`
	Type         FieldType `json:"type"`
	Length       int       `json:"length"`
	IsPrimaryKey bool      `json:"isPrimaryKey"`
	Nullish      bool      `json:"nullish"` // Equivalent to nullable
	AutoGen      bool      `json:"autoGen"`
	Default      any       `json:"default"`
	Attributes   string    `json:"attributes"`
	// Conditional children
}

// Index defines database indexes for an entity
// multiple columns can be separated by coma
type Index struct {
	Name      string `json:"name"`
	Columns   string `json:"column"`
	Unique    *bool  `json:"unique"`
	IndexType string `json:"indexType"`
}

// Relationship defines external key constraints between entities
type ForeignKeys struct {
	Name            string            `json:"name"`
	LocalField      string            `json:"localField"`
	ReferencedTable string            `json:"referencedTable"`
	ReferencedField string            `json:"referencedField"`
	OnDelete        ReferentialAction `json:"onDelete"`
	OnUpdate        ReferentialAction `json:"onUpdate"`
}

type ReferencedColumns struct {
	ConstraintName string            `json:"constraintName"` // Name of the foreign key constraint
	LocalColumn    string            `json:"localColumn"`    // Column in current table that has the FK
	ForeignTable   string            `json:"foreignTable"`   // Which table this column references
	ForeignColumn  string            `json:"foreignColumn"`  // Which column in the foreign table
	OnDelete       ReferentialAction `json:"onDelete"`       // CASCADE, SET NULL, RESTRICT, etc.
	OnUpdate       ReferentialAction `json:"onUpdate"`       // CASCADE, SET NULL, RESTRICT, etc.
}

// SourceTarget specifies an entity and a field within it
type SourceTarget struct {
	Entity string `json:"entity"`
	Field  string `json:"field"`
}
