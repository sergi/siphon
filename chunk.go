package p2b

type Chunk struct {
	//ID        xid.ID    `json:"id"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Host      string `json:"host"`
}
