package datasource

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/quexer/utee"
)

type Opt struct {
	MySqlConn string
	Debug     bool
}

type Ds struct {
	Db *gorm.DB
}

func gormDb() *gorm.DB {
	db, err := gorm.Open("mysql", "root:12345678@(localhost:3306)/stu?charset=utf8mb4&parseTime=True&loc=Local")
	utee.Chk(err)
	db.DB().SetMaxIdleConns(500)
	db.DB().SetMaxOpenConns(1500)
	db.SingularTable(true)
	db.LogMode(true)
	return db
}

func CreateDs() *Ds {
	ds := &Ds{
		Db: gormDb(),
	}
	fmt.Println("init gorm db")
	initAndMigration(ds)
	fmt.Println("init migration gorm db")
	return ds
}

func NotFound(err error) bool {
	return gorm.ErrRecordNotFound == err
}

func Duplicate(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}

func RunTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
