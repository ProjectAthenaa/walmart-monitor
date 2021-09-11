package main

import (
	"github.com/ProjectAthenaa/sonic-core/protos/monitor"
	monitor2 "github.com/ProjectAthenaa/walmart-monitor/monitor"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)



func main(){
	var lis net.Listener

	if os.Getenv("DEBUG") == "1"{
		lis, _ = net.Listen("tcp", ":4000")
	}else
	{
		lis, _ = net.Listen("tcp", ":3000")
	}
	server := grpc.NewServer()

	monitor.RegisterMonitorServer(server, monitor2.Server{})
	if err := server.Serve(lis); err != nil{
		log.Fatal(err)
	}

}