package model

import (
	"time"
)

type Message struct {
	Number   int64
	Username string
	Data     string
	Time     time.Time
}
