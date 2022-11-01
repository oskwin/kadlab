package d7024e

import (
	"errors"
	"log"
)

func (node Node) RPCPing(ping Ping, response *Contact) error {
	ping.Target = &node.Kad.RT.me
	log.Println("adding the pinging target to the correct bucket....")
	node.Kad.RT.AddContact(*ping.Requester)
	log.Println("ping received, sending pong....")
	*response = *ping.Target
	return nil

	// log.Println("contact not found....")
	// return errors.New("requested ID not found") // TODO get the error msg to actually print
}

// finds the closest contacts for the target
func (node Node) RPCFindContacts(network *Network, response *[]Contact) error {
	node.Kad.RT.AddContact(*network.Requester)
	// find the responders closest nodes to the target and store it in a list
	log.Println("RPC Accepted from ", network.Requester)
	//node.AddedNodes = node.Kad.RT.FindClosestContacts(network.Target.ID, k)
	list := ContactCandidates{
		contacts: node.Kad.RT.FindClosestContacts(network.Target.ID, k),
	}
	// respond with a list of contacts
	*response = list.contacts
	return nil
}

// find and send the requested data from the RPC receiver
func (node *Node) RPCFindValue(key *string, response *File) error {
	log.Println("key accepted, serching storage for key-value pair")
	if node.Kad.DataStore.Data == nil {
		log.Println("storage empty....")
		return errors.New("storage empty")
	} else if node.Kad.DataStore.Data[*key] == nil {
		log.Println("empty slot in map, this node does not hold requested value....")
		return errors.New("empty slot in map, this node does not hold requested value....")
	} else {
		*response = *node.Kad.DataStore.Data[*key]
		return nil
	}
}

func (node *Node) RPCStoreData(data *File, response *File) error {
	log.Println("RPC accepted....\nstore data at given key....")
	if node.Kad.DataStore.Data == nil {
		node.Kad.DataStore.Data = map[string]*File{}
	}
	node.Kad.DataStore.Data[*data.Key] = data
	*response = *data
	return nil
}

func PrintContacts(list []Contact) {
	log.Println("this nodes closest contacts")
	for i := range list {
		log.Println(list[i].ID)
	}
}
