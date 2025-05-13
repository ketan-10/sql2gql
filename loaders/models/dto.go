package models

type TableDTO struct {
	TableName string
	Columns   []*Column
}

type TableWithIndex struct {
	Table   *TableDTO
	Indexes []*Index
}

type TableRelations struct {
	*TableWithIndex
	ForeignKeys          []*ForeignKey // foregin keys of this table, one to Many
	ForeignKeysRef       []*ForeignKey // other tables pointing to our table, Many To One
}

type EnumDTO struct {
	*Enum
	DatabaseName string
	Values       []string
}

type AllModels struct {
	Enums  []*EnumDTO
	Tables []*TableDTO
}


type RltsDTO struct {
	*TableRelations
	UniqueTablesForRepoDependency []string
}