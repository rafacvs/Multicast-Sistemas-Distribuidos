package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "../pkg"
)

func main() {
	id := flag.Int("id", 0, "ID do processo é obrigatório.")
	configFile := flag.String("config", "peers.json", "Path para o arquivo de configuracao de peers")
	flag.Parse()

	if *id == 0 {
		log.Fatal("ID do processo (--id) é obrigatório.")
	}

	fmt.Printf("Iniciando processo (ID %d)\n", *id)
	fmt.Printf("Config: %s\n", *configFile)

	// TODO: Implementar logica

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("[EXECUTANDO] Aperte Ctrl+C para sair.")
	<-sigChan
	fmt.Println("\n [FINALIZADO].")
}
