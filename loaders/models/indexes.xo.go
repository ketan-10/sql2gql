package models

type IndexColumn struct {
	ColumnName string
	SequenceNo int
	Column     *Column // for ease of reach
}

type Index struct {
	IndexName string
	IsUnique  bool
	TableName string
	Columns   []*IndexColumn
}

type indexTable struct {
	IndexName  string
	IsUnique   bool
	ColumnName string
	SequenceNo int
}

func MySqlIndexes(db IDB, databaseName string, tableName string) ([]*Index, error) {

	// just load indexes, so if a index have multiple columns only single row wil be returned
	// Unique will let us know if we have multiple columns or not.
	const sqlstr = `SELECT ` +
		`index_name, ` +
		`NOT non_unique AS is_unique, ` +
		`column_name, ` +
		`seq_in_index AS seq_no ` +
		`FROM information_schema.statistics ` +
		`WHERE index_schema = ? AND table_name = ? `
		// + `AND index_name <> 'PRIMARY'`
		// Why we adding above row?

	q, err := db.Query(sqlstr, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	resultConverter := make(map[string][]indexTable)

	for q.Next() {

		var indexName, columnName string
		var seqNo int
		var isUnique bool

		// scan
		err = q.Scan(&indexName, &isUnique, &columnName, &seqNo)
		if err != nil {
			return nil, err
		}

		resultConverter[indexName] = append(resultConverter[indexName], indexTable{
			ColumnName: columnName,
			SequenceNo: seqNo,
			IndexName:  indexName,
			IsUnique:   isUnique,
		})

	}

	res := []*Index{}
	for indexName, values := range resultConverter {

		i := Index{}
		i.TableName = tableName
		i.IndexName = indexName

		// it will never be empty.
		i.IsUnique = values[0].IsUnique

		columns := []*IndexColumn{}

		for _, col := range values {
			columns = append(columns, &IndexColumn{
				ColumnName: col.ColumnName,
				SequenceNo: col.SequenceNo,
			})
		}
		i.Columns = columns

		res = append(res, &i)
	}

	return res, nil

}
