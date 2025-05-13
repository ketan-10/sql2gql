package loaders

import (
	_ "github.com/go-sql-driver/mysql" // empty import, to load drivers
	"github.com/ketan-10/mysql2graphql/internal"
	"github.com/ketan-10/mysql2graphql/loaders/models"
)

// init is a special function like main, which will be auto-called on start
func init() {

	// Register mysql loader into the system.
	internal.AllLoaders[internal.MYSQL] = &internal.LoaderImp{
		EnumList: models.MysqlEnums,
		DatabaseName: models.MysqlDatabaseName,
		EnumValueList: models.MysqlEnumValueList,
		TableList: models.MySqlTables,
		ColumList: models.MySqlColumns,
		IndexList: models.MySqlIndexes,
		ForeignKeysList: models.MySqlForeignKeys,
	}
}
