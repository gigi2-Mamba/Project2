package startup

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"project0/pkg/loggerDefine"
)

func InitDB() *gorm.DB {
	// dsn 很烦人
	//root:lxj360@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", "root", "root", "localhost", 13316, "webook")
	//log.Println(dsn)
	//gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Option... 结构， 无限填充
	Db, err := gorm.Open(mysql.Open("root:rootlxj0@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True"), &gorm.Config{})

	//log.Println("dsn: ",dsn)
	if err != nil {
		fmt.Println("db连接失败，error= ", err.Error())
		panic(err)
	}
	return Db

}
