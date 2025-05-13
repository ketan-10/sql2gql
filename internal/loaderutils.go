package internal

import (
	"github.com/ketan-10/mysql2graphql/loaders/models"
	"github.com/ketan-10/mysql2graphql/utils"
)

// --> input
// {
// 	name: "hello"
// 	id: 1,
// 	unique: true,
// 	columns: [
// 		{name: col1, no: 1},
// 		{name: col1, no: 2},
// 		{name: col1, no: 3},
// 	]
// }

// -- output
// {
// 	name: "hello"
// 	id: 1,
// 	unique: true,
// 	columns: [
// 		{name: col1, no: 1},
// 		{name: col1, no: 2},
// 		{name: col1, no: 3},
// 	]
// }
// {
// 	name: "hello"
// 	id: 1,
// 	unique: true,
// 	columns: [
// 		{name: col1, no: 1},
// 	]
// }
// {
// 	name: "hello"
// 	id: 1,
// 	unique: true,
// 	columns: [
// 		{name: col1, no: 1},
// 		{name: col1, no: 2},
// 	]
// }

// A multi-column index can still be effective even if you are only searching by a single column that is part of index.
// For example if you create an index on column (A, B, C). Mysql will create 3 seperate index (A), (A, B), (A, B, C)

func expandIndex(indexes []*models.Index) []*models.Index {
	var all_indexes []*models.Index
	for _, index := range indexes {
		all_indexes = append(all_indexes, index)

		for i := 1; i < len(index.Columns); i++ {
			new_idx := *index
			new_idx.IsUnique = false
			new_idx.Columns = utils.Filter[*models.IndexColumn](index.Columns, func(v *models.IndexColumn) bool {
				return v.SequenceNo <= i
			})
			all_indexes = append(all_indexes, &new_idx)
		}

	}

	return filterIndexes(all_indexes)
}

// {{ $tableNameCamel }}By{{range .Columns}}{{camelCase .ColumnName}}{{end}}
// Filter due to duplicate function names
func filterIndexes(indexes []*models.Index) []*models.Index {
	var filtered_indexes []*models.Index
	unique := make(map[string]*models.Index)
	for _, index := range indexes {
		fullColumns := ""
		for _, col := range index.Columns {
			fullColumns += col.ColumnName
		}
		unique[fullColumns] = index
	}
	for _, v := range unique {
		filtered_indexes = append(filtered_indexes, v)
	}
	return filtered_indexes
}

func filterIndexesOnlyFirstColumn(indexes []*models.Index) []*models.Index {
	var filtered_indexes []*models.Index
	unique := make(map[string]*models.Index)
	for _, index := range indexes {
		column := index.Columns[0].ColumnName
		unique[column] = index
	}
	for _, v := range unique {
		filtered_indexes = append(filtered_indexes, v)
	}
	return filtered_indexes
}

// attach more detials for ease of use in template
func attachColumnDetailsToIndex(indexes []*models.Index, table *models.TableDTO) {
	for _, index := range indexes {
		for _, col := range index.Columns {
			for _, tableCol := range table.Columns {
				if tableCol.ColumnName == col.ColumnName {
					col.Column = tableCol
				}
			}
		}
	}
}

// attach more detials for ease of use in template
func attachDetailsToForeignKeys(foreignKeys []*models.ForeignKey, table *models.TableDTO, tableWithIndexes []*models.TableWithIndex) {
	for _, key := range foreignKeys {
		key.Table = table
		for _, tableCol := range table.Columns {
			if tableCol.ColumnName == key.ColumnName {
				key.Column = tableCol
			}
		}
	}

	for _, key := range foreignKeys {
		for _, tableAndIndex := range tableWithIndexes {
			if tableAndIndex.Table.TableName == key.RefTableName {
				key.RefTable = tableAndIndex.Table

				for _, column := range tableAndIndex.Table.Columns {
					if column.ColumnName == key.RefColumnName {
						key.RefColumn = column
					}
				}

			}

		}
	}
}

func attachManyToOneForeignKeys(res []*models.TableRelations) {
	for _, tableRelation := range res {
		for _, key := range tableRelation.ForeignKeys {
			for _, tr := range res {
				if tr.Table.TableName == key.RefTableName {
					tr.ForeignKeysRef = append(tr.ForeignKeysRef, key)
				}
			}
		}
	}
}

func getUniqueRepoDependeciesTableNameForRLTS(tableRelation *models.TableRelations) []string {

	mymap := make(map[string]bool)
	var res []string
	// This Table Pointing to Other Table!!!, ManyToOne <- As Many records from other table can point to this table one record
	for _, fk := range tableRelation.ForeignKeys {
		if !mymap[fk.RefTableName] {
			mymap[fk.RefTableName] = true
			res = append(res, fk.RefTableName)
		}
	}

	//  Other Table Pointing to This Table!!!, OneToMany <- As This Table record can point to Multiple Other table record
	for _, fk := range tableRelation.ForeignKeysRef {
		if !mymap[fk.Table.TableName] {
			mymap[fk.Table.TableName] = true
			res = append(res, fk.Table.TableName)
		}
	}
	return res
}
