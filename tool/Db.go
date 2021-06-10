package tool

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GlobalMysqlDb *gorm.DB

type Db struct {

}


func (d *Db)connectionMysql() error {
	config := new(Config)
	mysqlConf ,err := config.GetSection("mysql")

	dsn := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		mysqlConf["root"],
		mysqlConf["pwd"],
		mysqlConf["host"],
		mysqlConf["port"],
		mysqlConf["db_trademark"],
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("数据库连接错误")
	}

	GlobalMysqlDb = db
	return nil
}


func (d *Db) GetDb() (*gorm.DB,error) {
	if GlobalMysqlDb == nil{
		err := d.connectionMysql()
		if err != nil {
			return nil,err
		}
	}
	return  GlobalMysqlDb,nil
}
