package siphon

// Chunk is a data structure that represents the chunks of data transferred by the client
// via UDP
type Chunk struct {
	ID        string `json:"id"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Host      string `json:"host"`
}
