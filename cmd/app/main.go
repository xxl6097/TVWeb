package main

import (
	"github.com/xxl6097/TVWeb/cmd/app/service"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs"
	"go.uber.org/zap"

	// Import buffer to register ukey for service signing
	_ "github.com/xxl6097/go-service/assets/buffer"
)

func main() {
	err := gs.Run(&service.TVWebService{})
	z.L().Debug("程序结束", zap.Error(err))
}
