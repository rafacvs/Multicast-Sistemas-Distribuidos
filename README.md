## Sistemas Distribuídos - Multicast Ordenado

### Proposta

Implementar um sistema de mensagens em grupo onde processos enviam e recebem mensagens via multicast. Cada mensagem carrega um timestamp e só pode ser entregue à aplicação quando:

- Estiver no início de uma fila ordenada
- For reconhecida (ACK) por todos os processos do grupo

O sistema deve suportar o cenário em que um ACK pode chegar antes da mensagem original em alguns processos, mantendo a consistência da ordem de entrega entre todos. O objetivo é demonstrar um multicast totalmente ordenado.

### Como executar

1. Configurar `peers.json`. Exemplo:

```json
[
  { "id": 1, "addr": "127.0.0.1:9001" },
  { "id": 2, "addr": "127.0.0.1:9002" },
  { "id": 3, "addr": "127.0.0.1:9003" }
]
```

2. Execute cada processo em um terminal separado:

```bash
# Terminal 1
go run cmd/main.go --id=1 --config=peers.json

# Terminal 2
go run cmd/main.go --id=2 --config=peers.json

# Terminal 3
go run cmd/main.go --id=3 --config=peers.json
```

### Objetivo

- Implementar um algoritmo de multicast totalmente ordenado.
- Garantir que todas as mensagens sejam entregues na mesma ordem em todos os processos.

### Aprendizado Envolvido

- Conceitos de relógios lógicos (Lamport) e ordenação total de eventos em sistemas distribuídos.
- Coordenação entre processos independentes por meio de mensagens (multicast conceitual).
- Noções de consistência na entrega de mensagens e condições necessárias para entrega segura (ordem + ACKs).
- Entendimento do desafio de mensagens fora de ordem (ex.: ACK antes da DATA) e como lidar com este caso.
