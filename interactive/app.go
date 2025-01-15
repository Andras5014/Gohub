package main

import (
	"github.com/Andras5014/webook/pkg/grpcx"
	"github.com/Andras5014/webook/pkg/saramax"
)

type App struct {
	Server    *grpcx.Server
	Consumers []saramax.Consumer
}
