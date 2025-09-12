package pkg

type MessageID struct {
	SenderID int   `json:"SenderID"` // ID do processo que enviou a mensagem
	Seq      int64 `json:"Seq"`      // Número de sequência local do remetente
}

// Mensagem na rede
type Envelope struct {
	Type      string    `json:"type"`      // "DATA" ou "ACK"
	FromID    int       `json:"fromId"`    // ID remetente
	Timestamp int64     `json:"timestamp"` // Lamport
	ID        MessageID `json:"id"`
	Payload   string    `json:"payload,omitempty"` // payload apenas para DATA
}

type Peer struct {
	ID   int    `json:"id"`
	Addr string `json:"addr"` // IP:porta
}

type MsgState struct {
	ID            MessageID
	DataTimestamp int64
	Payload       *string // Conteúdo da mensagem | nil
	Enqueued      bool
	AckedBy       map[int]struct{}
}

type QueueItem struct {
	ID            MessageID
	DataTimestamp int64
}
