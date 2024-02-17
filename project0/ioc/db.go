package ioc

import (
	"fmt"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
	"project0/pkg/gormx"
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
	//log.Println(" can out put db ",d.Dsn)
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

	// GORM自带的把控连接
	Db.Use(prometheus.New(prometheus.Config{
		DBName: "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				//一般都没什么卵用
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "society_pay",
		Subsystem: "webook",
		Name:      "gorm_db",
		Help:      "统计 GORM 的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	err = Db.Use(cb)
	if err != nil {
		panic(err)
	}

	return Db
}

// 函数延伸类型实现接口
type gormLoggerFunc func(string, ...loggerDefine.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...interface{})  {
	 g(msg,loggerDefine.Field{Key: "args",Val: fields})
}
