package pkg

import (
	"fmt"
)

type Node struct {
	ID       int
	Peers    []Peer
	Clock    *LamportClock
	Queue    *HoldBackQueue
	States   map[MessageID]*MsgState
	NextSeq  int64
	NumPeers int
}

func NewNode(ID int, Peers []Peer) *Node {
	Clock := NewLamportClock()
	Queue := NewHoldBackQueue()
	States := make(map[MessageID]*MsgState)
	NextSeq := int64(1)

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

func (n *Node) OnSendApp(payload string) {
	n.Clock.Tick()

	ID := MessageID{
		SenderID: n.ID,
		Seq:      n.NextSeq,
	}

	Envelope := Envelope{
		Type:      "DATA",
		FromID:    n.ID,
		Timestamp: n.Clock.Now(),
		ID:        ID,
		Payload:   payload,
	}

	fmt.Printf("(%d) Enviando: %+v\n", n.ID, Envelope)

	for i := 0; i < n.NumPeers; i++ {
		// TODO implementar envio para peers
	}

}

func (n *Node) OnReceiveDATA(env Envelope) {
	fmt.Printf("(%d) [DATA]: %+v\n", n.ID, env)
}

func (n *Node) OnReceiveACK(env Envelope) {
	fmt.Printf("(%d) [ACK]: %+v\n", n.ID, env)
}

func tryDeliver(n *Node) {}
