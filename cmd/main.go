package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"sd/t1/pkg"
)

func LoadPeers(configFile string) ([]pkg.Peer, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	var peers []pkg.Peer
	if err := json.Unmarshal(data, &peers); err != nil {
		return nil, fmt.Errorf("erro ao processar JSON de peers: %v", err)
	}

	for i, p := range peers {
		if p.ID <= 0 {
			return nil, fmt.Errorf("peer na posição %d tem ID inválido: %d", i, p.ID)
		}
		if p.Addr == "" {
			return nil, fmt.Errorf("peer na posição %d tem endereço vazio", i)
		}
	}

	return peers, nil
}

func main() {
	id := flag.Int("id", 0, "ID do processo é obrigatório.")

	delayDataMean := flag.Int("delay-data-mean-ms", 0, "Atraso médio (ms) para mensagens DATA")
	delayDataJitter := flag.Int("delay-data-jitter-ms", 0, "Jitter (±ms) aplicado ao atraso de DATA")
	delayAckMean := flag.Int("delay-ack-mean-ms", 0, "Atraso médio (ms) para mensagens ACK")
	delayAckJitter := flag.Int("delay-ack-jitter-ms", 0, "Jitter (±ms) aplicado ao atraso de ACK")
	randSeed := flag.Int64("rand-seed", 0, "Semente para aleatoriedade reprodutível (0 desabilita)")

	flag.Parse()

	if *id == 0 {
		log.Fatal("ID do processo (--id) é obrigatório.")
	}

	fmt.Printf("Iniciando processo (ID %d)\n", *id)
	peers, err := LoadPeers("peers.json")
	if err != nil {
		log.Fatalf("Falha ao carregar peers: %v", err)
	}

	fmt.Printf("Configuração carregada: %d peers encontrados\n", len(peers))
	for _, p := range peers {
		fmt.Printf("  - Peer %d: %s\n", p.ID, p.Addr)
	}

	node := pkg.NewNode(*id, peers)
	fmt.Printf("Node inicializado com ID %d\n", node.ID)

	pkg.SetNetworkDelays(*delayDataMean, *delayDataJitter, *delayAckMean, *delayAckJitter, *randSeed)

	err = node.SetupNetwork()
	if err != nil {
		log.Fatalf("Erro ao configurar rede: %v", err)
	}

	fmt.Println("[EXECUTANDO] Digite 'send <mensagem>' para enviar uma mensagem, ou Ctrl+C para sair.")

	inputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Erro ao ler stdin: %v\n", err)
				continue
			}

			text = strings.TrimSpace(text)
			inputChan <- text
		}
	}()

	fmt.Print("> ")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	running := true
	for running {
		select {
		case input := <-inputChan:
			handleCommand(input, node)
			fmt.Print("> ")
		case <-sigChan:
			fmt.Println("\n [FINALIZANDO]...")
			running = false
		}
	}
	fmt.Println("[FINALIZADO].")
}

func handleCommand(input string, node *pkg.Node) {
	parts := strings.SplitN(input, " ", 2)

	if len(parts) == 0 {
		return
	}

	cmd := parts[0]

	switch cmd {
	case "send":
		if len(parts) < 2 {
			fmt.Println("Uso: send <mensagem>")
			return
		}
		payload := parts[1]
		node.OnSendApp(payload)

	case "queue":
		fmt.Printf("Fila: %d mensagens\n", node.Queue.Size())

		items := node.Queue.Items()
		for index, item := range items {
			fmt.Printf("\t(%d): ID=%+v ts=%d\n", index+1, item.ID, item.DataTimestamp)
			fmt.Printf("\t > %s\n", *node.States[item.ID].Payload)
		}

	default:
		fmt.Printf("Comando desconhecido: %s\n", cmd)
		fmt.Println("Comandos disponíveis: send, queue")
	}
}
