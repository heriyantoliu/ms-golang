module accountservice

go 1.14

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.15

require (
	github.com/boltdb/bolt v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/h2non/gock v0.0.0-00010101000000-000000000000
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
)
