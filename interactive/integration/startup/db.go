package startup

import (
	"github.com/Andras5014/webook/config"
	"github.com/Andras5014/webook/interactive/repository/dao"

	"github.com/Andras5014/webook/pkg/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config, l logx.Logger) *gorm.DB {

	db, err := gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{
		//Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
		//	LogLevel:                  glogger.Info,
		//	IgnoreRecordNotFoundError: true,
		//	SlowThreshold:             time.Millisecond * 10,
		//}),
	})
	//err = db.Use(gormprometheus.New(gormprometheus.Config{
	//	DBName:          "webook",
	//	RefreshInterval: 10,
	//	StartServer:     false,
	//	MetricsCollector: []gormprometheus.MetricsCollector{
	//		&gormprometheus.MySQL{
	//			VariableNames: []string{"Threads_running"},
	//		},
	//	},
	//}))
	//
	//// 监控查询执行时间
	//callback := gormx.NewCallbacks(prometheus.SummaryOpts{
	//
	//	Namespace: "andras",
	//	Subsystem: "webook",
	//	Name:      "gorm_db",
	//	Help:      "GORM query duration in milliseconds",
	//	ConstLabels: map[string]string{
	//		"instance_id": "gorm_db_instance",
	//	},
	//	Objectives: map[float64]float64{
	//		0.5:   0.01,
	//		0.75:  0.01,
	//		0.9:   0.01,
	//		0.99:  0.001,
	//		0.999: 0.0001,
	//	},
	//})
	//
	//err = db.Use(callback)
	//if err != nil {
	//	panic(err)
	//}
	//db.Use(tracing.NewPlugin(tracing.WithDBName("echohub"),
	//	tracing.WithQueryFormatter(func(query string) string {
	//		l.Debug("query", logx.Any("query", query))
	//		return query
	//	}),
	//	// 不记录metrics和查询参数
	//	tracing.WithoutMetrics(),
	//	tracing.WithoutQueryVariables()))

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logx.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...interface{}) {
	g(msg, logx.Field{Key: "fields", Value: fields})
}
