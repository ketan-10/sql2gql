package models

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/ketan-10/mysql2graphql/utils"
)

type Column struct {
	FieldIndex   int            // field_ordinal
	ColumnName   string         // column_name
	DataType     string         // data_type (if enum ? ColumnName : realDataType)
	RealDataType string         // mySql data type eg int(11), varchar(255), enum(yes, no, maybe)
	Extra        string         // Extra info like (auto_increment, STORED GENERATED)
	IsEnum       bool           // is enumn
	NotNullable  bool           // is nullable?
	DefaultValue sql.NullString // default_value
	IsPrimaryKey bool           // is_primary_key
	IsGenerated  bool           // is auto generated

	// custom
	GoType      string
	GraphQLType string
	TableName   string
}

func MySqlColumns(db IDB, databaseName string, tableName string) ([]*Column, error) {
	var err error

	const sqlstr = `SELECT ` +
		`ordinal_position AS field_index, ` +
		`column_name, ` +
		`IF(data_type = 'enum', column_name, column_type) AS data_type, ` +
		`column_type AS real_data_type, ` +
		`extra, ` +
		`data_type = 'enum' AS is_enum, ` +
		`IF(is_nullable = 'YES', false, true) AS not_null, ` +
		`column_default AS default_value, ` +
		`IF(column_key = 'PRI', true, false) AS is_primary_key, ` +
		`IF(INSTR(EXTRA, 'GENERATED') > 0, true, false) AS is_generated ` +
		`FROM information_schema.columns ` +
		`WHERE table_schema = ? AND table_name = ? ` +
		`ORDER BY ordinal_position`

	// run query
	q, err := db.Query(sqlstr, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Column{}
	for q.Next() {
		c := Column{}
		c.TableName = tableName
		// scan
		err = q.Scan(&c.FieldIndex, &c.ColumnName,
			&c.DataType, &c.RealDataType,
			&c.Extra, &c.IsEnum,
			&c.NotNullable, &c.DefaultValue,
			&c.IsPrimaryKey, &c.IsGenerated)
		if err != nil {
			return nil, err
		}

		c.GoType = parseSQL2GoType(c)
		c.GraphQLType = parseGoToGraphqlType(c.GoType)

		res = append(res, &c)
	}

	return res, nil
}

func parseSQL2GoType(c Column) string {
	if c.IsEnum {
		return "enum." + utils.SnakeToCamel(c.TableName+"_"+c.ColumnName)
	}
	// remove unsigned
	dt := strings.TrimSuffix(c.DataType, " unsigned")

	dt, precision, _ := parsePrecision(dt)

	var nilVal, typ string
	nullable := !c.NotNullable
	switch dt {
	case "bit":
		nilVal = "0"
		if precision == 1 {
			nilVal = "false"
			typ = "bool"
			if nullable {
				nilVal = "sql.NullBool{}"
				typ = "sql.NullBool"
			}
			break
		} else if precision <= 8 {
			typ = "uint8"
		} else if precision <= 16 {
			typ = "uint16"
		} else if precision <= 32 {
			typ = "uint32"
		} else {
			typ = "uint64"
		}
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "bool", "boolean":
		nilVal = "false"
		typ = "bool"
		if nullable {
			nilVal = "sql.NullBool{}"
			typ = "sql.NullBool"
		}

	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "json":
		nilVal = `""`
		typ = "string"
		if nullable {
			nilVal = "sql.NullString{}"
			typ = "sql.NullString"
		}

	case "tinyint":
		//people using tinyint(1) really want a bool
		if precision == 1 {
			nilVal = "false"
			typ = "bool"
			if nullable {
				nilVal = "sql.NullBool{}"
				typ = "sql.NullBool"
			}
			break
		}
		nilVal = "0"
		typ = "int8"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "smallint":
		nilVal = "0"
		typ = "int16"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "mediumint", "int", "integer":
		nilVal = "0"
		typ = "int"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "bigint":
		nilVal = "0"
		typ = "int64"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "float":
		nilVal = "0.0"
		typ = "float32"
		if nullable {
			nilVal = "sql.NullFloat64{}"
			typ = "sql.NullFloat64"
		}

	case "decimal", "double":
		nilVal = "0.0"
		typ = "float64"
		if nullable {
			nilVal = "sql.NullFloat64{}"
			typ = "sql.NullFloat64"
		}

	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		typ = "[]byte"

	case "timestamp", "datetime", "date":
		nilVal = "time.Time{}"
		typ = "time.Time"
		if nullable {
			nilVal = "sql.NullTime{}"
			typ = "sql.NullTime"
		}

	case "time":
		// time is not supported by the MySQL driver. Can parse the string to time.Time in the user code.
		typ = "string"
	case "point":
		typ = "geo.Point"
		nilVal = "geo.Point{}"
	}

	// TODO assign to default value maybe 🤔
	_ = strings.ToUpper(nilVal)

	return typ
}

var PrecScaleRE = regexp.MustCompile(`\(([0-9]+)(\s*,[0-9]+)?\)$`)

// ParsePrecision extracts (precision[,scale]) strings from a data type and
// returns the data type without the string.
func parsePrecision(dt string) (string, int, int) {
	var err error

	precision := -1
	scale := -1

	m := PrecScaleRE.FindStringSubmatchIndex(dt)
	if m != nil {
		// extract precision
		precision, err = strconv.Atoi(dt[m[2]:m[3]])
		if err != nil {
			panic("could not convert precision")
		}

		// extract scale
		if m[4] != -1 {
			scale, err = strconv.Atoi(dt[m[4]+1 : m[5]])
			if err != nil {
				panic("could not convert scale")
			}
		}

		// change dt
		dt = dt[:m[0]] + dt[m[1]:]
	}

	return dt, precision, scale
}

func parseGoToGraphqlType(typ string) string {
	switch typ {
	case "int":
		fallthrough
	case "uint":
		fallthrough
	case "int64":
		typ = "Int"
	case "string":
		typ = "String"
	case "int8":
		typ = "IntBool"
	case "float64":
		typ = "Float"
	case "bool":
		typ = "Boolean"
	case "time.Time":
		typ = "Datetime"
	case "sql.NullTime":
		typ = "NullTime"
	case "sql.NullInt64":
		typ = "NullInt64"
	case "sql.NullFloat64":
		typ = "NullFloat64"
	case "sql.NullString":
		typ = "NullString"
	case "sql.NullBool":
		typ = "NullBool"
	case "geo.Point":
		typ = "Point"
	}
	return typ
}
