package store

import (
	"database/sql"
	"github.com/alextilot/golang-htmx-chatapp/model"
	"log"
	"time"
)

type MessageStore struct {
	DB *sql.DB
}

func NewMessageStore(db *sql.DB) (*MessageStore, error) {
	store := &MessageStore{DB: db}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS messages (username text not null, content text, time INTEGER);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return store, err
	}

	return store, nil
}

func (ms *MessageStore) Create(msg *model.Message) error {
	sqlStmt := `
	INSERT INTO messages (username, content, time) VALUES (?, ?, ?)
	`
	_, err := ms.DB.Exec(sqlStmt, msg.Username, msg.Data, msg.Time.UnixMilli())
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (ms *MessageStore) GetMostRecent(count int) ([]*model.Message, error) {
	messages := []*model.Message{}

	rows, err := ms.DB.Query("SELECT * FROM messages ORDER BY time DESC LIMIT ?", count)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return messages, err
	}

	for rows.Next() {
		var username string
		var content string
		var msec int64

		err = rows.Scan(&username, &content, &msec)
		if err != nil {
			log.Println(err)
			return messages, err
		}

		messages = append(messages, &model.Message{
			Username: username,
			Data:     content,
			Time:     time.UnixMilli(msec),
		})
	}

	return messages, nil
}
