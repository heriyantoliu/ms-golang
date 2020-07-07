package main

import (
	"flag"

	"github.com/heriyantoliu/ms-golang/common/messaging"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/heriyantoliu/ms-golang/common/config"
	"github.com/heriyantoliu/ms-golang/accountservice/dbclient"
	"github.com/heriyantoliu/ms-golang/accountservice/service"
	"github.com/spf13/viper"
)

var appName = "accountservice"

func init() {

	profile := flag.String("profile", "dev", "Environment profile, something similar to spring profiles")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	configServerUrl := flag.String("configServerUrl", "http://localhost:8888", "Address to config server")
	configBranch := flag.String("configBranch", "P9", "git branch to fetch configuration from")
	flag.Parse()

	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)

}

func main() {
	logrus.Infof("Starting %v\n", appName)

	config.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))

	initializeBoltClient()
	initializeMessaging()
	handleSigTerm(func(){
		service.MessagingClient.Close()
	})

	service.StartWebServer(viper.GetString("server_port"))
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url") {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}

	service.MessagingClient = &messaging.MessagingClient{}
	service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
	service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}

func handleSigTerm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
