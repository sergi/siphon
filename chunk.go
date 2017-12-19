package p2b

import (
	"encoding/json"
	"time"
)

type Chunk struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
}

func (c *Chunk) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, rawStrings["timestamp"])
	if err != nil {
		return err
	}
	c.Timestamp = t

	return nil
}

func (c *Chunk) MarshalJSON() ([]byte, error) {
	basicChunk := struct {
		Data      string `json:"url"`
		Timestamp string `json:"timestamp"`
		Host      string `json:"host"`
	}{
		Data:      c.Data,
		Timestamp: c.Timestamp.Format(time.RFC3339),
		Host:      c.Host,
	}

	return json.Marshal(basicChunk)
}
