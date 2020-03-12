package pkg

import (
	"os"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/ut"
)

func defaultDsOpt() *datasource.Opt {
	opt := &datasource.Opt{
		Debug:     getEnv("DEBUG", true).(bool),
		MySqlConn: getEnv("MYSQLCONN", "root:12345678@(localhost:3306)/stu?charset=utf8mb4&parseTime=True&loc=Local").(string),
	}
	return opt
}

func defaultWebOpt() *ut.WebOpt {
	opt := &ut.WebOpt{
		Addr: getEnv("ADDR", ":3000").(string),
		War:  getEnv("WAR", "./public").(string),
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
