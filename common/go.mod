module github.com/heriyantoliu/ms-golang/common

go 1.14

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.15

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/eapache/go-resiliency v1.2.0
	github.com/h2non/gock v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.3.0
)
