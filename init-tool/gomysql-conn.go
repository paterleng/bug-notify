package init_tool

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"strconv"
)

func GoMysqlConn() (*canal.Canal, error) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = Conf.MySQLConfig.Host + ":" + strconv.Itoa(Conf.MySQLConfig.Port)
	cfg.User = Conf.MySQLConfig.User
	cfg.Password = Conf.MySQLConfig.Password
	cfg.Dump.TableDB = Conf.Table.TableDB
	cfg.Dump.Tables = Conf.Table.TableName
	var tables []string
	for _, name := range Conf.Table.TableName {
		tables = append(tables, Conf.Table.TableDB+"."+name)
	}
	cfg.IncludeTableRegex = tables

	c, err := canal.NewCanal(cfg)
	return c, err
}
