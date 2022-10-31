package main

import (
	"kadlab/d7024e"
	"net/rpc"
)

func main() {
	node := new(d7024e.Node)

	me := d7024e.NewContact(d7024e.NewKademliaID("8111112300000000000000000000000000000000"), "localhost:8002")
	rt := d7024e.NewRoutingTable(me)

	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("5122111200000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("9122111400000000000000000000000000000000"), "localhost:8003"))
	for i := 0; i <= 50; i++ {
		rt.AddContact(d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost:8001"))
	}
	node.Kad.RT = rt

	rpc.Register(node)
	rpc.HandleHTTP()
	d7024e.Listen("localhost", 8002)

	for {

	}
}
