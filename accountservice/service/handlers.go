package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/heriyantoliu/ms-golang/accountservice/dbclient"
	"github.com/heriyantoliu/ms-golang/accountservice/model"
	"github.com/heriyantoliu/ms-golang/common/messaging"
)

var DBClient dbclient.IBoltClient
var MessagingClient messaging.IMessagingClient
var isHealthy = true
var client = &http.Client{}

func init() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	var accountId = mux.Vars(r)["accountId"]

	account, err := DBClient.QueryAccount(accountId)
	account.ServedBy = getIP()

	quote, err := getQuote()
	if err == nil {
		account.Quote = quote
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	notifVIP(account)

	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func notifVIP(account model.Account) {
	if account.Id == "10000" {
		go func(account model.Account) {
			vipNotification := model.VipNotification{AccountId: account.Id, ReadAt: time.Now().UTC().String()}
			data,  _ := json.Marshal(vipNotification)
			logrus.Infof("Notifiying VIP account %v\n", account.Id)
			err := MessagingClient.PublishOnQueue(data, "vip_queue")
			if err != nil {
				logrus.Infoln(err.Error())
			}
		}(account)
	}
}

func getQuote() (model.Quote, error) {
	req, _ := http.NewRequest("GET", "http://quotes-service:8080/api/quote?strength=4", nil)
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == 200 {
		quote := model.Quote{}
		bytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bytes, &quote)
		return quote, nil
	} else {
		return model.Quote{}, fmt.Errorf("Some error")
	}
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbUp := DBClient.Check()
	if dbUp && isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccesible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	}
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])
	if err != nil {
		logrus.Infoln("Invalid request to SetHealthyState, allowed values are true or false")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}
