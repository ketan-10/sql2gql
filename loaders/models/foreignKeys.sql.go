package models

type ForeignKey struct {
	ForeignKeyName string
	ColumnName     string
	RefTableName   string
	RefColumnName  string

	// added after
	Column    *Column
	RefColumn *Column
	Table     *TableDTO
	RefTable  *TableDTO
}

func MySqlForeignKeys(db IDB, databaseName string, tableName string) ([]*ForeignKey, error) {
	var err error

	const sqlstr = `SELECT ` +
		`constraint_name AS foreign_key_name, ` +
		`column_name AS column_name, ` +
		`referenced_table_name AS ref_table_name, ` +
		`referenced_column_name AS ref_column_name ` +
		`FROM information_schema.key_column_usage ` +
		`WHERE referenced_table_name IS NOT NULL AND table_schema = ? AND table_name = ?`

	// run query
	q, err := db.Query(sqlstr, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ForeignKey{}
	for q.Next() {
		fk := ForeignKey{}

		// scan
		err = q.Scan(&fk.ForeignKeyName, &fk.ColumnName, &fk.RefTableName, &fk.RefColumnName)
		if err != nil {
			return nil, err
		}

		res = append(res, &fk)
	}

	return res, nil
}
