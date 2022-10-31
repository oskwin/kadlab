package d7024e

type Node struct {
	Kad        Kademlia
	AddedNodes []Contact
}

func RemoveMe(m Kademlia, c []Contact) []Contact {
	for i := len(c) - 1; i >= 0; i-- {
		//log.Println(c[i].ID, m.RT.me.ID)
		if c[i].ID.Equals(m.RT.me.ID) {
			c = append(c[:i], c[i+1:]...)
		}
	}
	return c
}
