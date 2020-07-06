package service

import (
	"github.com/sirupsen/logrus"

	"net/http"
)

func StartWebServer(port string) {

	r:= NewRouter()
	http.Handle("/", r)
	logrus.Infoln("Starting HTTP service at " + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		logrus.Infoln("An error occured starting HTTP listener at port " + port)
		logrus.Infoln("Error: " + err.Error())
	}
}
