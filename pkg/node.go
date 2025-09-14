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
	Network  *UDPNetwork
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
	timestamp := n.Clock.Tick()

	ID := MessageID{
		SenderID: n.ID,
		Seq:      n.NextSeq,
	}
	n.NextSeq++

	envelope := Envelope{
		Type:      "DATA",
		FromID:    n.ID,
		Timestamp: timestamp,
		ID:        ID,
		Payload:   payload,
	}

	fmt.Printf("(%d) Enviando mensagem: %+v\n", n.ID, envelope)

	if n.Network != nil {
		n.Network.Broadcast(envelope)
	} else {
		fmt.Printf("(%d) AVISO: Network não configurado, mensagem não enviada\n", n.ID)
	}
}

func (n *Node) OnReceiveDATA(env Envelope) {
	n.Clock.OnReceive(env.Timestamp)

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

	fmt.Printf("(%d) Enviando ACK para mensagem: ID=%+v ts=%d\n", n.ID, env.ID, tsAck)

	if n.Network != nil {
		n.Network.Broadcast(ackEnv)
	} else {
		fmt.Printf("(%d) AVISO: Network não configurado, ACK não enviado\n", n.ID)
	}
}

func (n *Node) OnReceiveACK(env Envelope) {
	n.Clock.OnReceive(env.Timestamp)

	st, exists := n.States[env.ID]
	if !exists {
		st = &MsgState{
			ID:      env.ID,
			AckedBy: make(map[int]struct{}),
		}
		n.States[env.ID] = st
	}
	st.AckedBy[env.FromID] = struct{}{}

	n.tryDeliver()
}

func (n *Node) SetupNetwork() error {
	var localAddr string
	for _, peer := range n.Peers {
		if peer.ID == n.ID {
			localAddr = peer.Addr
			break
		}
	}

	if localAddr == "" {
		return fmt.Errorf("endereço local não encontrado para ID %d", n.ID)
	}

	network, err := NewUDPNetwork(localAddr, n)
	if err != nil {
		return fmt.Errorf("erro ao configurar rede: %v", err)
	}

	n.Network = network
	network.StartReceiving()

	fmt.Printf("(%d) Rede configurada no endereço %s\n", n.ID, localAddr)
	return nil
}

func (n *Node) tryDeliver() {
	for !n.Queue.IsEmpty() {
		top := n.Queue.Peek()
		st, exists := n.States[top.ID]

		if exists && len(st.AckedBy) == n.NumPeers {
			msg := n.Queue.Pop()

			fmt.Printf("(%d) [ENTREGUE] ID=%+v ts=%d payload='%s' ACKs=%d/%d\n",
				n.ID, msg.ID, msg.DataTimestamp, *st.Payload,
				len(st.AckedBy), n.NumPeers)

			delete(n.States, msg.ID)

			continue
		} else {
			if exists {
				fmt.Printf("(%d) [AGUARDANDO] ID=%+v ACKs=%d/%d\n", n.ID, top.ID, len(st.AckedBy), n.NumPeers)
			} else {
				fmt.Printf("(%d) [ERRO] Estado não encontrado para ID=%+v\n", n.ID, top.ID)
			}
			break
		}
	}
}
