package init_tool

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	*ProjectConfig `mapstructure:"project"`
	*LogConfig     `mapstructure:"log"`
	*MySQLConfig   `mapstructure:"mysql"`
	*Table         `mapstructure:"table"`
}

type Table struct {
	TableDB   string   `mapstructure:"table_db"`
	TableName []string `mapstructure:"table_name"`
}

type ProjectConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Address   string `mapstructure:"address"`
	Port      string `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         string `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// 初始化viper，用于解析配置文件
func ViperInit() error {
	viper.SetConfigFile("./config/config.yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("修改了配置文件...")
		viper.Unmarshal(&Conf)
	})
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("readconfig failed,err: %v", err)
		return err
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		fmt.Errorf("unmarshal to Conf failed, err: %v", err)
		return err
	}
	return err
}
