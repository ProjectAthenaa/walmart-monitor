package monitor

import (
	"context"
	"github.com/ProjectAthenaa/sonic-core/protos/monitor"
	monitor_controller "github.com/ProjectAthenaa/sonic-core/protos/monitorController"
	proxy_rater "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/product"
	"google.golang.org/grpc"
	"os"
)

var (
	rdb = core.Base.GetRedis("cache")
	pxClient, _ = sonic.NewPerimeterXClient("localhost:3000")
	proxyClient proxy_rater.ProxyRaterClient
)

type Server struct {
	monitor.UnimplementedMonitorServer
}

func init(){
	if os.Getenv("DEBUG") == "1"{
		conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
		if err != nil{
			panic(err)
		}

		proxyClient = proxy_rater.NewProxyRaterClient(conn)
		return
	}

	conn, err := grpc.Dial("proxy-rater.general.svc.cluster.local:3000", grpc.WithInsecure())
	if err != nil{
		panic(err)
	}
	proxyClient = proxy_rater.NewProxyRaterClient(conn)
}

func (s Server) Start(ctx context.Context, task *monitor_controller.Task) (*monitor_controller.BoolResponse, error) {
	t, err := NewTask(task)
	if err != nil {
		return nil, err
	}

	t.Start(product.SiteWalmart,proxyClient)

	return &monitor_controller.BoolResponse{
		Stopped: false,
		Error:   nil,
	}, nil
}
