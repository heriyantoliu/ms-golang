package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/heriyantoliu/ms-golang/common/config"
	"github.com/heriyantoliu/ms-golang/common/messaging"
	"github.com/heriyantoliu/ms-golang/vipservice/service"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
)

var appName = "vipservice"
var messagingClient messaging.IMessagingClient

func init() {
	profile := flag.String("profile", "dev", "Environment profile, something similar to spring profiles")
	logrus.SetFormatter(&logrus.JSONFormatter{})
	configServerUrl := flag.String("configServerUrl", "http://localhost:8888", "address to config server")
	configBranch := flag.String("configBranch", "P9", "git branch to fetch configuration from")
	flag.Parse()

	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)
}

func main() {
	logrus.Infoln("Starting " + appName + "...")

	config.LoadConfigurationFromBranch(viper.GetString("configServerUrl"), appName, viper.GetString("profile"), viper.GetString("configBranch"))
	initializeMessaging()

	handleSigTerm(func() {
		if messagingClient != nil {
			messagingClient.Close()
		}
	})
	service.StartWebServer(viper.GetString("server_port"))
}

func onMessage(delivery amqp.Delivery) {
	logrus.Infof("Got a message: %v\n", string(delivery.Body))
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url" ){
		panic("No 'broker_url set in configuratio, cannot start")
	}
	messagingClient = &messaging.MessagingClient{}
	messagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))

	err := messagingClient.SubscribeToQueue("vip_queue", appName, onMessage)
	failOnError(err, "Could not start subsribe to vip_queue")

	err = messagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
	failOnError(err, "Could not start subscribe to "+viper.GetString("config_event_bus")+" topic")
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

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Infof("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}