package sqlrender

type TableSchema struct {
	Catalog  string
	Name     string
	Location string
}

type TableRef struct {
	TableSchema TableSchema
	TableName   string
}

type Column struct {
	ColName string
	ColType string
	Comment string
}

type TableProperties struct {
	Format        string
	PartitionedBy []string
	Extra         map[string]any
}

type Table struct {
	TableSchema TableSchema
	TableName   string
	Columns     []Column
	TableProp   TableProperties
}

type AddColumnsData struct {
	Table   TableRef
	Columns []Column
}

type DropColumnsData struct {
	Table    TableRef
	ColNames []string
}

type CreateSchemaData struct {
	TableSchema TableSchema
	IfNotExists bool
}

type DropSchemaData struct {
	TableSchema TableSchema
	IfExists    bool
	Cascade     bool
}

type CreateTableData struct {
	Table       Table
	IfNotExists bool
}

type DropTableData struct {
	Table    TableRef
	IfExists bool
}

type RenameColumnData struct {
	Table   TableRef
	OldName string
	NewName string
}

type RenameTableData struct {
	Table        TableRef
	NewTableName string
}

type SetFileFormatData struct {
	Table  TableRef
	Format string
}

type SetPartitioningData struct {
	Table         TableRef
	PartitionedBy []string
}

type SetTableLocationData struct {
	Table            TableRef
	ExternalLocation string
}
