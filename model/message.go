package model

import (
	"time"
)

type Message struct {
	Username string
	Data     string
	Time     time.Time
}
