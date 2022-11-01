package d7024e

import "log"

const (
	k     = 20
	alpha = 3
)

// Kademlia is the node
type Kademlia struct {
	RT           *RoutingTable
	DataStore    DataStorage
	Candidates   ContactCandidates
	Net          Network
	ClosestNodes ContactCandidates
}

func (kademlia *Kademlia) LookupContact(contact *Contact) {
	kademlia.Net.Target = contact
	kademlia.Net.Requester = &kademlia.RT.me
	kademlia.RT.me.distance = kademlia.RT.me.ID.CalcDistance(contact.ID) // since the distance is not yet decided for me

	log.Println("getting my candidates closest to the target....")
	kademlia.Candidates.contacts = kademlia.RT.FindClosestContacts(contact.ID, alpha)
	//kademlia.Candidates.contacts = RemoveMe(*kademlia, kademlia.Candidates.contacts)

	kademlia.Lookup(contact.ID)
}

func (kademlia *Kademlia) LookupData(hash *string) {
	kademlia.Net.FindValue = nil
	log.Println("performing a lookup for the k-closest nodes of the given hash....")
	kademlia.Net.Requester = &kademlia.RT.me
	candidates := kademlia.RT.GetClosest(*hash) // gives the closest nodes in ContactList
	for i := range candidates {
		kademlia.Net.Target = &candidates[i]
		// target := candidates[i].ID.String()
		kademlia.Net.SendFindDataMessage(hash)
		if &kademlia.Net.FindValue != nil {
			if kademlia.DataStore.Data == nil {
				kademlia.DataStore.Data = map[string]*File{}
			}
			kademlia.DataStore.Data[*hash] = kademlia.Net.FindValue
			break

		}
	}
}

func (kademlia *Kademlia) Store(value *[]byte) {
	log.Println("generating key to value....")
	kademlia.Net.Key = NewRandomKademliaID().String()
	log.Println("performing a lookup to store a key-value pair in k-closest nodes to generated key....")
	kademlia.Lookup(NewKademliaID(kademlia.Net.Key))
	for i := range kademlia.Net.ContactList {
		kademlia.Net.Target = &kademlia.Net.ContactList[i]
		kademlia.Net.SendStoreMessage(value)
	}
	log.Println("key-value pair stored in k-closest nodes to the hash key....")
}

func (rt *RoutingTable) GetClosest(obj string) []Contact {
	return rt.FindClosestContacts(NewKademliaID(obj), k)
}

func (kademlia *Kademlia) Lookup(hash *KademliaID) {
	var contactedNodes ContactCandidates
	var (
		exitFlag    bool
		exitCounter int
	)
	var kClosestNodes struct {
		nodes  []Contact
		status [k]string
	}

	for !exitFlag {
		exitFlag = false
		for i := range kademlia.Candidates.contacts {
			kademlia.ClosestNodes.Append(kademlia.Net.SendFindContactMessage(&kademlia.Candidates.contacts[i]))
			for j := range kademlia.ClosestNodes.contacts {
				kademlia.ClosestNodes.contacts[j].distance = kademlia.ClosestNodes.contacts[i].ID.CalcDistance(hash)
			}
		}
		contactedNodes.Append(kademlia.Candidates.contacts)
		for i := range kademlia.ClosestNodes.contacts {
			kademlia.RT.AddContact(kademlia.ClosestNodes.contacts[i])
		}
		kClosestNodes.nodes = kademlia.RT.FindClosestContacts(hash, k)
		//kClosestNodes.nodes = RemoveMe(*kademlia, kClosestNodes.nodes)
		for i := range kademlia.Candidates.contacts {
			for j := range kClosestNodes.nodes {
				if kademlia.Candidates.contacts[i].ID == kClosestNodes.nodes[j].ID {
					kClosestNodes.status[j] = "contacted"
					exitCounter++
				}
			}
		}
		if exitCounter == len(kClosestNodes.nodes) {
			exitFlag = true // true
		} else {
			kademlia.Candidates.contacts = kClosestNodes.nodes
			exitFlag = false // false
		}
		// for i := range kClosestNodes.nodes {
		// 	log.Println(kClosestNodes.nodes[i], kClosestNodes.status[i])
		// }
		exitCounter = 0
	}
	kademlia.Net.ContactList = append(kClosestNodes.nodes)
}
