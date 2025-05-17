package internal

import (
	"fmt"
	"strings"

	"github.com/ketan-10/sql2gql/loaders/models"
	"github.com/ketan-10/sql2gql/templates"
)

// DatabaseType enum currently only mysql is supported.
type DatabaseType uint16

const MYSQL = DatabaseType(1)

func (lt *DatabaseType) String() string {
	switch *lt {
	case MYSQL:
		return "mysql"
	}
	return ""
}

// The loader interface
type ILoader interface {
	LoadSchema(*Args) error
}

var AllLoaders = map[DatabaseType]ILoader{}

// The loader implementation
// Drivers like mysql will create object for this.
type LoaderImp struct {
	EnumList        func(db models.IDB, databaseName string) ([]*models.Enum, error)
	DatabaseName    func(db models.IDB) (string, error)
	EnumValueList   func(db models.IDB, databaseName string, tableName string, columnName string) (string, error)
	TableList       func(db models.IDB, databaseName string) ([]string, error)
	ColumList       func(db models.IDB, databaseName string, tableName string) ([]*models.Column, error)
	IndexList       func(db models.IDB, databaseName string, tableName string) ([]*models.Index, error)
	ForeignKeysList func(db models.IDB, databaseName string, tableName string) ([]*models.ForeignKey, error)
}

// Entry point to load everything
func (lt *LoaderImp) LoadSchema(args *Args) error {
	var err error

	database, err := lt.loadDatabaseName(args)
	if err != nil {
		return err
	}
	args.DatabaseName = database

	allEnumDTO, err := lt.loadEnums(args)
	if err != nil {
		return err
	}

	tables, err := lt.loadTables(args)
	if err != nil {
		return err
	}

	tableWithIndexes, err := lt.loadIndex(args, tables)
	if err != nil {
		return err
	}

	tableRelations, err := lt.loadForeignKeys(args, tableWithIndexes)
	if err != nil {
		return err
	}

	allModels := &models.AllModels{
		Enums:  allEnumDTO,
		Tables: tables,
	}

	// execute enum
	for _, enum := range allEnumDTO {
		err := args.ExecuteTemplate(templates.ENUM, fmt.Sprintf("%s_%s", enum.TableName, enum.ColumnName), enum)
		if err != nil {
			return err
		}
	}

	// execute Table
	for _, tableWithIndex := range tableWithIndexes {
		tableDto := &models.TableWithIndex{
			Table:   tableWithIndex.Table,
			Indexes: filterIndexesOnlyFirstColumn(tableWithIndex.Indexes),
		}
		err := args.ExecuteTemplate(templates.TABLE, tableWithIndex.Table.TableName, tableDto)
		if err != nil {
			return err
		}
	}

	// execute repos
	for _, tableWithIndex := range tableWithIndexes {
		err = args.ExecuteTemplate(templates.REPO, tableWithIndex.Table.TableName+"_repository", tableWithIndex)
		if err != nil {
			return err
		}
	}

	// execute repos
	err = args.ExecuteTemplate(templates.QUERY_MUTATION, "query_mutation", tableWithIndexes)
	if err != nil {
		return err
	}
	
	// execute rlts
	for _, tableRelation := range tableRelations {
		uniqueTableNames := getUniqueRepoDependeciesTableNameForRLTS(tableRelation)
		rltsDto := &models.RltsDTO{
			TableRelations:                tableRelation,
			UniqueTablesForRepoDependency: uniqueTableNames,
		}
		err = args.ExecuteTemplate(templates.RLTS, tableRelation.Table.TableName+"_rlts_repository", rltsDto)
		if err != nil {
			return err
		}
	}

	// execute graphql schema
	for _, tableRelation := range tableRelations {

		err = args.ExecuteTemplate(templates.GRAPH_SCHEMA, tableRelation.Table.TableName, tableRelation)
		if err != nil {
			return err
		}
	}

	// execute go_resolver.go
	err = args.ExecuteTemplate(templates.GO_RESOLVER, "go_resolver", tableRelations)
	if err != nil {
		return err
	}

	// execute wire.go
	err = args.ExecuteTemplate(templates.GO_WIRE, "wire.gen", tableRelations)
	if err != nil {
		return err
	}

	// execute scalar.graphql
	err = args.ExecuteTemplate(templates.ENUM_SCALAR, "scalar", allModels)
	if err != nil {
		return err
	}

	// execute gqlgen.yml
	err = args.ExecuteTemplate(templates.GQLGEN, "gqlgen", allModels)
	if err != nil {
		return err
	}

	return nil
}

func (lt *LoaderImp) loadDatabaseName(args *Args) (string, error) {
	if lt.DatabaseName == nil {
		return "", fmt.Errorf("schema name loader is not implemented for %s", args.DatabaseType.String())
	}
	return lt.DatabaseName(args.DB)
}

func (lt *LoaderImp) loadIndex(args *Args, tables []*models.TableDTO) ([]*models.TableWithIndex, error) {

	res := []*models.TableWithIndex{}

	for _, table := range tables {
		indexes, err := lt.IndexList(args.DB, args.DatabaseName, table.TableName)

		all_indexes := expandIndex(indexes)
		// add column details to index for ease of use
		attachColumnDetailsToIndex(all_indexes, table)

		if err != nil {
			return nil, err
		}

		tableWithIndex := &models.TableWithIndex{
			Table:   table,
			Indexes: all_indexes,
		}

		res = append(res, tableWithIndex)
	}
	return res, nil
}

func (lt *LoaderImp) loadForeignKeys(args *Args, tablesAndIndexes []*models.TableWithIndex) ([]*models.TableRelations, error) {

	res := []*models.TableRelations{}

	for _, tablesAndIndex := range tablesAndIndexes {

		// type ForeignKey struct {
		// 	ForeignKeyName string
		// 	ColumnName     string
		// 	RefTableName   string
		// 	RefColumnName  string
		//
		//  Column Column
		//  RefColumn Column
		//  RefTable Table
		// }

		foreignKeys, err := lt.ForeignKeysList(args.DB, args.DatabaseName, tablesAndIndex.Table.TableName)
		if err != nil {
			return nil, err
		}

		// add column details to foreign keys for ease of use
		attachDetailsToForeignKeys(foreignKeys, tablesAndIndex.Table, tablesAndIndexes)

		// empty ref
		var foreignKeyRef []*models.ForeignKey

		tableRelations := &models.TableRelations{
			TableWithIndex: tablesAndIndex,
			ForeignKeys:    foreignKeys,
			ForeignKeysRef: foreignKeyRef,
		}

		res = append(res, tableRelations)
	}

	attachManyToOneForeignKeys(res)

	return res, nil
}

func (lt *LoaderImp) loadTables(args *Args) ([]*models.TableDTO, error) {
	tables, err := lt.TableList(args.DB, args.DatabaseName)
	if err != nil {
		return nil, err
	}
	var allTableDTO []*models.TableDTO

	for _, table := range tables {
		columns, err := lt.ColumList(args.DB, args.DatabaseName, table)
		if err != nil {
			return nil, err
		}
		allTableDTO = append(allTableDTO, &models.TableDTO{
			TableName: table,
			Columns:   columns,
		})
	}
	return allTableDTO, nil
}

func (lt *LoaderImp) loadEnums(args *Args) ([]*models.EnumDTO, error) {
	enums, err := lt.EnumList(args.DB, args.DatabaseName)
	if err != nil {
		return nil, err
	}

	var allEnumDTO []*models.EnumDTO
	for _, e := range enums {
		// fmt.Printf("%s, %s \n", e.ColumnName, e.TableName)
		enumValues, err := lt.loadEnumValues(args, e)

		if err != nil {
			return nil, err
		}

		allEnumDTO = append(allEnumDTO, &models.EnumDTO{
			Enum:         e,
			DatabaseName: args.DatabaseName,
			Values:       enumValues,
		})
	}
	return allEnumDTO, nil
}

func (lt *LoaderImp) loadEnumValues(args *Args, enum *models.Enum) ([]string, error) {
	if lt.EnumValueList == nil {
		return nil, fmt.Errorf("enumValue loader is not implemented for %s", args.DatabaseType.String())
	}

	values, err := lt.EnumValueList(args.DB, args.DatabaseName, enum.TableName, enum.ColumnName)

	if err != nil {
		return nil, err
	}
	// value is in 'A','B','C' we want to convert to a list
	list := strings.Split(values[1:len(values)-1], "','")
	return list, nil
}
