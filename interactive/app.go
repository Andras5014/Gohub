package main

import (
	"github.com/Andras5014/gohub/pkg/grpcx"
	"github.com/Andras5014/gohub/pkg/saramax"
)

type App struct {
	Server    *grpcx.Server
	Consumers []saramax.Consumer
}
