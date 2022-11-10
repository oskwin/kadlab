package main

import (
	"fmt"
	"kadlab/d7024e"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// set some prerequisites
	var input string
	var flag bool
	host := GetOutboundIP().String()
	node := new(d7024e.Node)

	me := d7024e.NewContact(d7024e.NewRandomKademliaID(), GetOutboundIP().String()+":4000")
	rt := d7024e.NewRoutingTable(me)
	node.Kad.RT = rt
	node.Kad.DataStore.Data = make(map[string]*d7024e.File)
	myNet := new(d7024e.Network)
	myNet.Requester = &me

	node = &d7024e.Node{
		Kad: d7024e.Kademlia{
			RT:  rt,
			Net: *myNet,
		},
	}
	log.Println("welcome to node ", host, me.ID)

	log.Println("adding the bootstrapnode....")
	node.Kad.RT.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "130.240.77.15:4000"))

	rpc.Register(node)
	rpc.HandleHTTP()
	d7024e.Listen(host, 4000)

	// send a lookup request of a random node with the IP that this node is hosted on
	log.Println("sending a random node id to the bootstrap to join the network....")
	node.Kad.LookupContact(&me)
	log.Println(node.Kad.Net.ContactList)
	log.Println("network joined.... pinging the contacts....")
	// the RPC responds with a list of the closest contacts that the RPC has
	for i := range node.Kad.Net.ContactList {
		err := node.Kad.Net.SendPingMessage(&node.Kad.Net.ContactList[i])
		if err != nil {
			log.Println(err)
		}
		log.Println("The pinged node is alive! Adding it to the bucketlist....")
		for i := range node.Kad.Net.ContactList {
			node.Kad.RT.AddContact(node.Kad.Net.ContactList[i])
		}
	}

	flag = false
	for !flag {
		// command line interface - [put] [get] [exit]
		log.Println("choose a command to excecute....\n[put.....]\n[get.....]\n[exit....]\n[list....]")
		fmt.Scan(&input)

		switch {

		case input == "put":
			//var hash string
			var msg string
			node.Kad.Net.Requester = &me
			log.Println(me.ID)
			fmt.Println("enter a message that you want to store")
			fmt.Scan(&msg)
			t := []byte(msg)
			node.Kad.Store(&t)
		case input == "get":
			var hash string
			fmt.Println("enter a hash key of the value you want to fetch....")
			fmt.Scan(&hash)
			node.Kad.LookupData(&hash)
			log.Println(node.Kad.DataStore)
			log.Println("=================================================================================================")
		case input == "exit":
			flag = true
			log.Println("shutting down node....")
			log.Println("=================================================================================================")
		case input == "list":
			d7024e.PrintContacts(node.Kad.Net.ContactList)
			for key, value := range node.Kad.DataStore.Data {
				log.Println(key, value)
			}
		}
	}
}

func fetchTarget(id string, list []d7024e.Contact) d7024e.Contact {
	var contactFound d7024e.Contact
	for i := range list {
		if id == list[i].ID.String() {
			contactFound = list[i]
			return contactFound
		}
	}
	return contactFound
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
