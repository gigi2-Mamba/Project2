package mysqlInit

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

/*
orm框架，连接数据库
*/

func Init() (err error) {
	// dsn 很烦人
	//root:lxj360@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
	//172.31.227.133
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", "root", "root", "localhost", 13316, "webook")
	//log.Println(dsn)
	gorm.Open(mysql.Open(dsn), &gorm.Config{})

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("连接失败，error= ", err.Error())
		return
	}
	return

}

func Close() {
	db, _ := Db.DB()

	db.Close()
}
