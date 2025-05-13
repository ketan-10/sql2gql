// Package models contains the types for schema 'public'.
package models

// Enum represents a enum.
type Enum struct {
	ColumnName string
	TableName  string
}

// Mysql Enums runs a custom query, returning results as Enum.
func MysqlEnums(db IDB, schema string) ([]*Enum, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`DISTINCT column_name AS enum_name, table_name ` +
		`FROM information_schema.columns ` +
		`WHERE data_type = 'enum' AND table_schema = ?`

	// run query
	q, err := db.Query(sqlstr, schema)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Enum{}
	for q.Next() {
		e := Enum{}

		// scan
		err = q.Scan(&e.ColumnName, &e.TableName)
		if err != nil {
			return nil, err
		}

		res = append(res, &e)
	}

	return res, nil
}
