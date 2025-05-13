package models

// returns the database name
func MysqlDatabaseName(db IDB) (string, error) {
	row := db.QueryRow("SELECT SCHEMA()")
	var data string
	err := row.Scan(&data)
	if err != nil {
		return "", err
	}
	return data, nil
}
