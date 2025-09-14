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
	configFile := flag.String("config", "peers.json", "Path para o arquivo de configuracao de peers")
	flag.Parse()

	if *id == 0 {
		log.Fatal("ID do processo (--id) é obrigatório.")
	}

	fmt.Printf("Iniciando processo (ID %d)\n", *id)
	fmt.Printf("Config: %s\n", *configFile)

	peers, err := LoadPeers(*configFile)
	if err != nil {
		log.Fatalf("Falha ao carregar peers: %v", err)
	}

	fmt.Printf("Configuração carregada: %d peers encontrados\n", len(peers))
	for _, p := range peers {
		fmt.Printf("  - Peer %d: %s\n", p.ID, p.Addr)
	}

	node := pkg.NewNode(*id, peers)
	fmt.Printf("Node inicializado com ID %d\n", node.ID)

	err = node.SetupNetwork()
	if err != nil {
		log.Fatalf("Erro ao configurar rede: %v", err)
	}

	inputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Erro ao ler stdin: %v\n", err)
				continue
			}

			text = strings.TrimSpace(text)
			inputChan <- text
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("[EXECUTANDO] Digite 'send <mensagem>' para enviar uma mensagem, ou Ctrl+C para sair.")

	running := true
	for running {
		select {
		case input := <-inputChan:
			handleCommand(input, node)
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
		fmt.Printf("Enviando: %s\n", payload)
		node.OnSendApp(payload)

	case "queue":
		fmt.Printf("Fila: %d mensagens\n", node.Queue.Size())
		if !node.Queue.IsEmpty() {
			top := node.Queue.Peek()
			st, exists := node.States[top.ID]
			if exists {
				fmt.Printf("Topo: ID=%+v ts=%d ACKs=%d/%d\n",
					top.ID, top.DataTimestamp, len(st.AckedBy), node.NumPeers)
			}
		}

	default:
		fmt.Printf("Comando desconhecido: %s\n", cmd)
		fmt.Println("Comandos disponíveis: send, queue")
	}
}
