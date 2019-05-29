package main

import (
	"gt-monitor/cmd/conv-server/internal"
	"gt-monitor/common/zap"
	"gt-monitor/config"
	"gt-monitor/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	pflag.String("conf", "./config.yaml", "set configuration `file`")
	pflag.String("profile", "dev", "app profile")
	pflag.String("log-dir", "./logs", "server logs dir")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	configFile := viper.GetString("conf")
	var logger = zap.Get()
	if utils.IsEmpty(configFile) {
		pflag.Usage()
		logger.Info("can't not found config file")
		os.Exit(1)
	}

	cnf, err := config.ReadConfig(configFile)
	if err != nil {
		logger.Error("read cnf file error. file=", configFile)
		log.Fatalf("read cnf file error: %s", err)
	}

	//
	server := internal.NewConvServer(cnf)
	if err := server.Init(); err != nil {
		logger.Error("ConvServer init error!", err)
		log.Fatalf("ConvServer init error: %s", err)
	}

	logger.Info("HTTP Server running ")

	server.Run()
}
