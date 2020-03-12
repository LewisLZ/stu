package pkg

import (
	"os"

	"liuyu/stu/pkg/datasource"
)

func defaultDsOpt() datasource.Opt {
	opt := datasource.Opt{
		Debug:     getEnv("debug", true).(bool),
		MySqlConn: getEnv("mysqlconn", "root:12345678@(localhost:3306)/stu?charset=utf8mb4&parseTime=True&loc=Local").(string),
	}
	return opt
}

func getEnv(key string, def interface{}) interface{} {
	env, have := os.LookupEnv(key)
	if !have {
		return def
	}
	return env
}
