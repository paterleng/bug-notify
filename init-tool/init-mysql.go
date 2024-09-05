package init_tool

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MysqlInit() (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", Conf.MySQLConfig.User, Conf.MySQLConfig.Password, Conf.MySQLConfig.Host, Conf.MySQLConfig.Port, Conf.MySQLConfig.DB)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	DB, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return
	} else {
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxOpenConns(Conf.MySQLConfig.MaxOpenConns)
		sqlDB.SetMaxIdleConns(Conf.MySQLConfig.MaxIdleConns)
	}
	return
}
