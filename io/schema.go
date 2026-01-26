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
	Entities      []Entity       `json:"entities"`
	Relationships []Relationship `json:"relationships"`
}

// Entity represents a database table/collection
type Entity struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	PrimaryKey  []string            `json:"primaryKey"`
	OnConflict  *OnConflictStrategy `json:"onConflict,omitempty"`
	Fields      []Field             `json:"fields"`
	Indexes     []Index             `json:"indexes,omitempty"`
}

// Field represents a column/field in an entity
type Field struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Subtype     *string   `json:"subtype,omitempty"`
	Length      *int      `json:"length,omitempty"`
	Nulish      bool      `json:"nulish"` // Equivalent to nullable
	Default     any       `json:"default,omitempty"`
	Embed       *bool     `json:"embed,omitempty"`
	Description string    `json:"description"`
	EnumValues  []string  `json:"enumValues,omitempty"`

	// Conditional children
	Fields    []Field `json:"fields,omitempty"`    // if type == object
	Items     *Field  `json:"items,omitempty"`     // if type == array
	ValueType *string `json:"valueType,omitempty"` // if type == map
}

// Index defines database indexes for an entity
type Index struct {
	Name        string   `json:"name"`
	Fields      []string `json:"fields"`
	Unique      *bool    `json:"unique,omitempty"`
	Sparse      *bool    `json:"sparse,omitempty"`
	Description string   `json:"description"`
}

// Relationship defines external key constraints between entities
type Relationship struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"` // e.g., "belongsTo"
	Source   SourceTarget      `json:"source"`
	Target   SourceTarget      `json:"target"`
	OnDelete ReferentialAction `json:"onDelete"`
	OnUpdate ReferentialAction `json:"onUpdate"`
}

// SourceTarget specifies an entity and a field within it
type SourceTarget struct {
	Entity string `json:"entity"`
	Field  string `json:"field"`
}
