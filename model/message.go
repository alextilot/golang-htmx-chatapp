package model

import (
	"time"
)

type Message struct {
	ClientID string
	Username string
	Data     string
	Time     time.Time
}
