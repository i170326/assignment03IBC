package assignment03IBC

import (
	a2 "assignment02IBC_i170326" //"github.com/i170326/assignment02IBC_i170326"
	"encoding/gob"
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
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}
			connList = append(connList, conn)
			HandleConnectionSatoshi(connList[totalConn])
			totalConn = totalConn + 1
		}
	}
	if user == "others" {
		ln, err := net.Listen("tcp", portno)
		if err != nil {
			log.Fatal(err)
		}
		for {
			peerNode, err = ln.Accept()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func HandleConnectionSatoshi(c net.Conn) {
	recvdSlice := make([]byte, 4096)
	c.Read(recvdSlice)
	neighboursPorts = append(neighboursPorts, string(recvdSlice))
	c.Write([]byte("Welcome to MishaalCoin. Kindly wait for other nodes to join."))
	if totalConn != 0 {
		chainHead = a2.InsertBlock("", "", "Satoshi", 0, chainHead)
	}
}

func WaitForQuorum() {
	for totalConn != Quorum {
	}

}

func WriteString(c net.Conn, myListeningAddress string) {
	c.Write([]byte(myListeningAddress))
	recvd := make([]byte, 4096)
	c.Read(recvd)
	fmt.Print(string(recvd))
}

func ReadString(c net.Conn) string {
	rcvbuff := make([]byte, 4096)
	c.Read(rcvbuff)
	return string(rcvbuff)
}
func SendChainandConnInfo() {
	fmt.Print("Quorum reached. Sending peer info to nodes")
	fmt.Print(neighboursPorts)
	for i := 0; i < totalConn; i++ {
		if i == totalConn-1 {
			connList[i].Write([]byte(neighboursPorts[0]))
		}
		if i != totalConn {
			connList[i].Write([]byte(neighboursPorts[i+1]))
		}
	}
	for i := 0; i < totalConn; i++ {
		gobEncoder := gob.NewEncoder(connList[i])
		err := gobEncoder.Encode(chainHead)
		if err != nil {
			log.Println(err)
		}
	}
}

func ReceiveChain(c net.Conn) *a2.Block {
	var rcvdBlock a2.Block
	dec := gob.NewDecoder(c)
	err := dec.Decode(&rcvdBlock)
	if err != nil {
		log.Println(err)
	}
	return &rcvdBlock
}
