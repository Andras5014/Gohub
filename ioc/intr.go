package ioc

import (
	interactivev1 "github.com/Andras5014/gohub/api/proto/gen/interactive/v1"
	"github.com/Andras5014/gohub/config"
	"github.com/Andras5014/gohub/interactive/service"
	"github.com/Andras5014/gohub/internal/web/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInteractiveGrpcClient(svc service.InteractiveService, cfg *config.Config) interactivev1.InteractiveServiceClient {
	var opts []grpc.DialOption
	if cfg.Grpc.Client.Intr.Secure {
		// 加载证书
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(cfg.Grpc.Client.Intr.Addr, opts...)
	if err != nil {
		panic(err)
	}
	local := client.NewInteractiveServiceAdapter(svc)
	remote := interactivev1.NewInteractiveServiceClient(cc)
	res := client.NewGreyScaleInteractiveServiceClient(local, remote, int32(cfg.Grpc.Client.Intr.Threshold))

	viper.OnConfigChange(func(in fsnotify.Event) {
		res.UpdateThreshold(cfg.Grpc.Client.Intr.Threshold)
	})
	return res
}
