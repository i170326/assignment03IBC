package assignment03IBC

import (
	a2 "assignment02IBC_i170326" //"github.com/i170326/assignment02IBC_i170326"
	"fmt"
	"log"
	"net"
)

var connList []net.Conn
var neighboursPorts []string
var totalConn = 0
var Quorum = 0
var SatoshiPort string
var peerNode net.Conn

var chainHead *a2.Block

func StartListening(portno string, user string) {
	if user == "satoshi" {
		chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
		SatoshiPort = portno
		ln, err := net.Listen("tcp", portno)
		if err != nil {
			log.Fatal(err)
		}
		connList[totalConn], err = ln.Accept()
		if err != nil {
			log.Println(err)
		}
		HandleConnectionSatoshi(connList[totalConn])
		totalConn = totalConn + 1
	}
	if user == "others" {
		ln, err := net.Listen("tcp", portno)
		if err != nil {
			log.Fatal(err)
		}
		peerNode, err = ln.Accept()
		if err != nil {
			log.Println(err)
		}
	}
}

func HandleConnectionSatoshi(c net.Conn) {
	recvdSlice := make([]byte, 11)
	c.Read(recvdSlice)
	neighboursPorts[totalConn] = string(recvdSlice)
	c.Write([]byte("Welcome to MishaalCoin. Kindly wait for other nodes to join."))
}

func WaitForQuorum() {
	for totalConn != Quorum {
		ln, err := net.Listen("tcp", SatoshiPort)
		if err != nil {
			log.Fatal(err)
		}
		connList[totalConn], err = ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		HandleConnectionSatoshi(connList[totalConn])
		totalConn = totalConn + 1
		chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
	}

}

func WriteString(c net.Conn, myListeningAddress string) {
	c.Write([]byte(myListeningAddress))
	recvd := make([]byte, 11)
	c.Read(recvd)
	fmt.Println(string(recvd))
}

func ReadString(c net.Conn) string {
	rcvbuff := make([]byte, 11)
	c.Read(rcvbuff)
	return string(rcvbuff)
}
func SendChainandConnInfo() {
	fmt.Println("Quorum reached. Sending peer info to nodes")
	for i := 0; i <= totalConn; i++ {
		if i == totalConn {
			connList[i].Write([]byte(neighboursPorts[0]))
		}
		if i != totalConn {
			connList[i].Write([]byte(neighboursPorts[i+1]))
		}
	}
}
