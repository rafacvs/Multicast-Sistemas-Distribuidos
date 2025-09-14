package pkg

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type UDPNetwork struct {
	LocalAddr *net.UDPAddr
	Conn      *net.UDPConn
	Node      *Node
}

var (
	delayDataMeanMs   int
	delayDataJitterMs int
	delayAckMeanMs    int
	delayAckJitterMs  int
	randomSource      *rand.Rand
)

func SetNetworkDelays(dataMeanMs, dataJitterMs, ackMeanMs, ackJitterMs int, seed int64) {
	delayDataMeanMs = dataMeanMs
	delayDataJitterMs = dataJitterMs
	delayAckMeanMs = ackMeanMs
	delayAckJitterMs = ackJitterMs
	if seed != 0 {
		randomSource = rand.New(rand.NewSource(seed))
	} else {
		randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

func sleepWithDelayFor(env Envelope) {
	var mean, jitter int
	if env.Type == "ACK" {
		mean, jitter = delayAckMeanMs, delayAckJitterMs
	} else {
		mean, jitter = delayDataMeanMs, delayDataJitterMs
	}
	if mean == 0 && jitter == 0 {
		return
	}
	// delay = mean ± jitter
	d := mean
	if jitter > 0 && randomSource != nil {
		delta := randomSource.Intn(2*jitter+1) - jitter
		d += delta
		if d < 0 {
			d = 0
		}
	}
	time.Sleep(time.Duration(d) * time.Millisecond)
}

func NewUDPNetwork(localAddr string, node *Node) (*UDPNetwork, error) {
	addr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("erro ao resolver endereço UDP: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir socket UDP: %v", err)
	}

	network := &UDPNetwork{
		LocalAddr: addr,
		Conn:      conn,
		Node:      node,
	}

	return network, nil
}

func (n *UDPNetwork) SendToPeer(peer Peer, env Envelope) error {
	sleepWithDelayFor(env)

	data, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("erro ao serializar envelope: %v", err)
	}

	peerAddr, err := net.ResolveUDPAddr("udp", peer.Addr)
	if err != nil {
		return fmt.Errorf("erro ao resolver endereço do peer: %v", err)
	}

	_, err = n.Conn.WriteToUDP(data, peerAddr)

	if err != nil {
		return fmt.Errorf("erro ao enviar para peer %d: %v", peer.ID, err)
	}

	return nil
}

func (n *UDPNetwork) Broadcast(env Envelope) {
	for _, peer := range n.Node.Peers {
		err := n.SendToPeer(peer, env)
		if err != nil {
			fmt.Printf("Erro ao enviar para peer %d: %v\n", peer.ID, err)
		}
	}
}

func (n *UDPNetwork) StartReceiving() {
	go func() {
		buffer := make([]byte, 4096)
		for {
			bytesRead, addr, err := n.Conn.ReadFromUDP(buffer)

			if err != nil {
				fmt.Printf("Erro ao ler do socket: %v\n", err)
				continue
			}

			var env Envelope
			err = json.Unmarshal(buffer[:bytesRead], &env)
			if err != nil {
				fmt.Printf("Erro ao desserializar mensagem de %s: %v\n", addr.String(), err)
				continue
			}

			n.handleMessage(env)
		}
	}()
	fmt.Println("Recebimento de mensagens iniciado...")
}

func (n *UDPNetwork) handleMessage(env Envelope) {
	fmt.Printf("\n[%s] Mensagem recebida de %d\n", env.Type, env.FromID)

	switch env.Type {
	case "DATA":
		n.Node.OnReceiveDATA(env)
	case "ACK":
		n.Node.OnReceiveACK(env)
	default:
		fmt.Printf("Tipo de mensagem desconhecido: %s\n", env.Type)
	}
}

func (n *UDPNetwork) Close() error {
	return n.Conn.Close()
}
