package startup

import (
	"github.com/Andras5014/gohub/pkg/logx"
	"go.uber.org/zap"
)

func InitLogger() logx.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logx.NewZapLogger(l)
}
