// Package portscanner defines an utility for scanning port on a target
package portscanner

import (
	"net"
	"strconv"
	"time"

	"github.com/alexiscampan/go.pkg/log"
	"go.uber.org/zap"
)

// OpenPorts function output the open ports on our target
func OpenPorts(ip string) error {
	log.Bg().Info("Scanning open ports for target", zap.String("target", ip))
	activeThreads := 0
	doneChannel := make(chan bool)

	for port := 0; port <= 65535; port++ {
		go testTCPConnection(ip, port, doneChannel)
		activeThreads++
	}
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
	return nil
}

func testTCPConnection(ip string, port int, doneChannel chan bool) {
	_, err := net.DialTimeout(
		"tcp", ip+":"+strconv.Itoa(port),
		time.Second*10,
	)
	if err == nil {
		log.Bg().Info("Port is open", zap.String("port", strconv.Itoa(port)))
	}
	doneChannel <- true
}
