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

	// TEST: on receive local
	n.OnReceiveDATA(Envelope)

	fmt.Printf("(%d) Node atualizado: %+v\n", n.ID, n)
}

func (n *Node) OnReceiveDATA(env Envelope) {
	ts := n.Clock.OnReceive(env.Timestamp)
	fmt.Printf("(%d) [DATA]: %+v | Clock: %d -> %d\n", n.ID, env, n.Clock.Now(), ts)

	st, exists := n.States[env.ID]
	if !exists {
		st = &MsgState{
			ID:      env.ID,
			AckedBy: make(map[int]struct{}),
		}
		n.States[env.ID] = st
	}

	st.DataTimestamp = env.Timestamp
	st.Payload = &env.Payload

	if !st.Enqueued {
		qi := &QueueItem{
			ID:            env.ID,
			DataTimestamp: env.Timestamp,
		}
		n.Queue.Push(qi)
		st.Enqueued = true
		fmt.Printf("(%d) Mensagem enfileirada: %+v\n", n.ID, qi)
	}

	tsAck := n.Clock.Tick()
	ackEnv := Envelope{
		Type:      "ACK",
		FromID:    n.ID,
		Timestamp: tsAck,
		ID:        env.ID,
	}
	fmt.Printf("(%d) Enviando ACK: %+v | Clock: %d -> %d\n", n.ID, ackEnv, n.Clock.Now(), tsAck)

	n.OnReceiveACK(ackEnv) // Processar local
	n.tryDeliver(&ackEnv)
}

func (n *Node) OnReceiveACK(env Envelope) {
	fmt.Printf("(%d) [ACK]: %+v\n", n.ID, env)
}

func (n *Node) tryDeliver(env *Envelope) {
	fmt.Printf("(%d) Tentando entregar mensagem: %+v\n", n.ID, env)
}
