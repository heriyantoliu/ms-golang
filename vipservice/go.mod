module github.com/heriyantoliu/ms-golang/vipservice

go 1.14

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.15

require (
	github.com/gorilla/mux v1.7.4
	github.com/heriyantoliu/ms-golang/common v0.0.0-20200707034000-b3a5a13e779e
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
)
