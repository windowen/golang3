package service

import (
	"serverApi/pkg/gozero/service"
	"serverApi/pkg/zlogger"
)

var group *service.ServiceGroup

func init() {
	group = service.NewServiceGroup()
}

func RegisterService(service service.Service) {
	group.Add(service)
}

func Start(info string) {
	zlogger.Infof("service start: %s", info)
	group.Start()
}

func Stop(info string) {
	group.Stop()
	zlogger.Infof("service stop: %s", info)
}
