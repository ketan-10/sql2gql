package models

// information_schema.column return "enum('A','B','C')"
func MysqlEnumValueList(db IDB, databaseName string, tableName string, columnName string) (string, error) {
	const query = `SELECT ` +
		`SUBSTRING(column_type, 6, CHAR_LENGTH(column_type) - 6) AS enum_values ` +
		`FROM information_schema.columns ` +
		`WHERE data_type = 'enum' AND table_schema = ? AND table_name = ? AND column_name = ?`

	var data string
	err := db.QueryRow(query, databaseName, tableName, columnName).Scan(&data)
	if err != nil {
		return "", err
	}
	return data, nil
}
