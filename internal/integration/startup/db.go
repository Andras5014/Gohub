package startup

import (
	"context"
	"fmt"
	"github.com/Andras5014/gohub/config"
	"github.com/Andras5014/gohub/internal/repository/dao"
	"github.com/Andras5014/gohub/pkg/logx"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config, l logx.Logger) *gorm.DB {

	db, err := gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			LogLevel:                  glogger.Warn,
			IgnoreRecordNotFoundError: true,
			SlowThreshold:             time.Millisecond * 10,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

var mongoDB *mongo.Database

func InitMongoDB() *mongo.Database {
	if mongoDB == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context,
				startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
		}
		opts := options.Client().
			ApplyURI("mongodb://root:root@localhost:27017").
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mongoDB = client.Database("gohub")
	}
	return mongoDB
}

type gormLoggerFunc func(msg string, fields ...logx.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...interface{}) {
	// 提取 GORM 日志中的关键信息
	if len(fields) >= 4 {
		// 格式化日志输出
		formattedMsg := fmt.Sprintf("[GORM] %.3fms | Rows: %v | SQL: %s", fields[2], fields[3], fields[4])
		g(formattedMsg, logx.Any("details", fields))
	} else {
		// 如果 fields 的数量少于 4 个，直接打印原始信息
		g(msg, logx.Any("args", fields))
	}
}
