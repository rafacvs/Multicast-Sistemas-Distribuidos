package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sd/t1/pkg"
)

// LoadPeers carrega a configuração de peers do arquivo JSON
func LoadPeers(configFile string) ([]pkg.Peer, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	var peers []pkg.Peer
	if err := json.Unmarshal(data, &peers); err != nil {
		return nil, fmt.Errorf("erro ao processar JSON de peers: %v", err)
	}

	// Validar que todos os peers têm ID e endereço
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

	// TODO: Implementar logica
	Node := pkg.NewNode(*id, peers)
	fmt.Printf("Node inicializado: %+v\n", Node)

	Node.OnSendApp("Hello Node")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("[EXECUTANDO] Aperte Ctrl+C para sair.")
	<-sigChan
	fmt.Println("\n [FINALIZADO].")
}
