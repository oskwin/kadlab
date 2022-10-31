package main

import (
	"kadlab/d7024e"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	node := new(d7024e.Node)
	node.Kad.DataStore.Data = make(map[string]d7024e.File)
	me := d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), GetOutboundIP().String()+":4000")
	rt := d7024e.NewRoutingTable(me)

	myNet := new(d7024e.Network)

	node.Kad.Net = *myNet
	node.Kad.RT = rt

	rpc.Register(node)
	rpc.HandleHTTP()
	d7024e.Listen(GetOutboundIP().String(), 4000)

	log.Println("this is the root node....")
	log.Println("Address: ", me.Address)
	for {

	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
