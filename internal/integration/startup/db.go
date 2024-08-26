package startup

import (
	"fmt"
	"github.com/Andras5014/webook/config"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config, l logger.Logger) *gorm.DB {

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

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, fields ...interface{}) {
	// 提取 GORM 日志中的关键信息
	if len(fields) >= 4 {
		// 格式化日志输出
		formattedMsg := fmt.Sprintf("[GORM] %.3fms | Rows: %v | SQL: %s", fields[2], fields[3], fields[4])
		g(formattedMsg, logger.Any("details", fields))
	} else {
		// 如果 fields 的数量少于 4 个，直接打印原始信息
		g(msg, logger.Any("args", fields))
	}
}
