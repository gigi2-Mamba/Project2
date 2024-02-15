package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"project0/pkg/loggerDefine"
)

func InitDB(l loggerDefine.LoggerV1) *gorm.DB {
	// dsn 很烦人
	//root:lxj360@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True"
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", "root", "root", "localhost", 13316, "webook")
	//log.Println(dsn)
	//gorm.Open(mysql.Open(dsn), &gorm.Config{})
	type dcfg struct {
		Dsn string `yaml:"dsn"`
	}
	var d dcfg
	viper.UnmarshalKey("db",&d,)
	log.Println(" xxxxx ",d.Dsn)
	// Option... 结构， 无限填充
	Db, err := gorm.Open(mysql.Open(d.Dsn), &gorm.Config{
		//Logger: glogger.New(gormLoggerFunc(l.Debug),glogger.Config{
		//	// 慢查询
		//	SlowThreshold: 0,
		//	LogLevel: glogger.Info,
		//},),

		})
	//Db.Callback().Query().Before("gorm:query").Register("my_plugin:sql_logger", func(db *gorm.DB) {
	//	fmt.Println(db.Statement.SQL)
	//})
	//Db.Set("CONNECT_TIMEOUT", 30*time.Second)
	if err != nil {
		fmt.Println("db连接失败，error= ", err.Error())
		panic(err)
	}
	return Db
}

// 函数延伸类型实现接口
type gormLoggerFunc func(string, ...loggerDefine.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...interface{})  {
	 g(msg,loggerDefine.Field{Key: "args",Val: fields})
}
