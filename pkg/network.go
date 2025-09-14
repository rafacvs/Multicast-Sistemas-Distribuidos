package pkg

import (
	"encoding/json"
	"fmt"
	"net"
)

type UDPNetwork struct {
	LocalAddr *net.UDPAddr
	Conn      *net.UDPConn
	Node      *Node
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

	fmt.Printf("Enviado para peer %d (%s): %v\n", peer.ID, peer.Addr, env.Type)
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

			go n.handleMessage(env)
		}
	}()
	fmt.Println("Recebimento de mensagens iniciado...")
}

func (n *UDPNetwork) handleMessage(env Envelope) {
	fmt.Printf("Mensagem recebida de %d: %s\n", env.FromID, env.Type)

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
