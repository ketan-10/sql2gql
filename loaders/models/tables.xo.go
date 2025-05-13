package models

func MySqlTables(db IDB, schema string) ([]string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`table_name ` +
		`FROM information_schema.tables ` +
		`WHERE table_schema = ?`

	// run query
	q, err := db.Query(sqlstr, schema)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []string
	for q.Next() {
		var tableName string

		// scan
		err = q.Scan(&tableName)
		if err != nil {
			return nil, err
		}

		res = append(res, tableName)
	}

	return res, nil
}
