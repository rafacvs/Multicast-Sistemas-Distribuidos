package pkg

type Node struct {
	ID       int
	Peers    []Peer
	Clock    *LamportClock
	Queue    *HoldBackQueue
	States   map[MessageID]*MsgState
	NextSeq  int
	NumPeers int
}

func NewNode(ID int, Peers []Peer) *Node {
	Clock := NewLamportClock()
	Queue := NewHoldBackQueue()
	States := make(map[MessageID]*MsgState)
	NextSeq := 1

	return &Node{
		ID:       ID,
		Peers:    Peers,
		Clock:    Clock,
		Queue:    Queue,
		States:   States,
		NextSeq:  NextSeq,
		NumPeers: len(Peers),
	}
}
