package status

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/message"
)

type PayloadNewStatus struct {
	ID        string          `json:"id"`
	Chat      string          `json:"chat"`
	Sender    string          `json:"sender"`
	Message   message.Message `json:"message"`
	Timestamp time.Time       `json:"timestamp"`
}
