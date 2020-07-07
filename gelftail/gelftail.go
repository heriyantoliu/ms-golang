package main

import (
	"encoding/json"
	"flag"
	"github.com/heriyantoliu/ms-golang/gelftail/aggregator"
	"github.com/heriyantoliu/ms-golang/gelftail/transformer"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"sync"
)

var authToken = ""
var port *string

func init() {
	data, err := ioutil.ReadFile("token.txt")
	if err != nil {
		msg := "Cannot find token.txt that should contain our Loggly token"
		logrus.Errorln(msg)
		panic(msg)
	}
	authToken = string(data)

	port = flag.String("port", "12202", "UDP port for the gelftail")

	flag.Parse()
}

func main() {
	logrus.Println("Starting Gelf-tail server...")
	ServerConn := startUDPServer(*port)
	defer ServerConn.Close()

	var bulkQueue = make(chan []byte, 1)

	go aggregator.Start(bulkQueue, authToken)
	go listenForLogStatements(ServerConn, bulkQueue)

	logrus.Infoln("Started Gelf-tail server")

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()


}

func startUDPServer(port string) *net.UDPConn {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":" +port)
	checkError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	checkError(err)

	return ServerConn
}

func checkError(err error) {
	if err != nil {
		logrus.Println("Error: ", err)
		os.Exit(0)
	}
}

func listenForLogStatements(ServerConn *net.UDPConn, bulkQueue chan[]byte) {
	buf := make([]byte, 8192)
	var item map[string]interface{}
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			logrus.Errorf("Problem reading UDP message into buffer: %v\n", err.Error())
			continue
		}

		err = json.Unmarshal(buf[0:n], &item)
		if err != nil {
			logrus.Errorln("Problem unmarshalling log message into JSON: " + err.Error())
			item = nil
			continue
		}



		processedLogMessage, err := transformer.ProcessLogStatement(item)
		if err != nil {
			logrus.Printf("Problem parsing message: %v", string(buf[0:n]))
		} else {
			bulkQueue <- processedLogMessage
		}
		item = nil
	}
}