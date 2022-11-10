package d7024e

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

type Network struct {
	Target      *Contact
	Requester   *Contact
	Recipient   *Contact
	ContactList []Contact
	FindValue   *File
	Key         string
}

type Ping Network
type Pong Network

// open up a server that keeps listening and serving requests
func Listen(ip string, port int) {
	ln, err := net.Listen("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
	}
	go http.Serve(ln, nil)
}

// sends a Ping to requested receiver
func (network *Network) SendPingMessage(contact *Contact) error {
	var responder *Contact
	network.Target = contact
	client, err := rpc.DialHTTP("tcp", contact.Address)
	if err != nil {
		log.Println(err)
	}
	err = client.Call("Node.RPCPing", &network, &responder)
	if err != nil {
		return err
	}
	return nil
}

// Dial and call a node to find the k closest nodes to the target contact
func (network *Network) SendFindContactMessage(contact *Contact) []Contact {
	var list *[]Contact

	// dial to the contactcandidate to see that nodes closest contacts to the target
	client, err := rpc.DialHTTP("tcp", contact.Address)
	if err != nil {
		log.Println("tcp connection error")
	}

	// make the call to the "server"
	client.Call("Node.RPCFindContacts", &network, &list)
	return *list
}

func (network *Network) SendFindDataMessage(hash *string) {
	var response File
	log.Println("sending message to find data to : ", network.Target.ID, network.Target.Address)
	log.Println(&hash)
	client, err := rpc.DialHTTP("tcp", network.Target.Address)
	if err != nil {
		log.Println(err)
	}
	client.Call("Node.RPCFindValue", &hash, &response)
	log.Println("value found: ", string(response.Value))
	if response.Value == nil {
		log.Println("empty value")
		return
	}
	network.FindValue = &response
}

func (network *Network) SendStoreMessage(value *[]byte) {
	var response File
	data := File{
		Key:   &network.Key, // Generated hash key
		Value: *value,
	}

	log.Println("sending message to store data to : ", network.Target.ID, network.Target.Address)

	client, err := rpc.DialHTTP("tcp", network.Target.Address)
	if err != nil {
		log.Println(err)
	}
	client.Call("Node.RPCStoreData", &data, &response)
	log.Println(response)
}
