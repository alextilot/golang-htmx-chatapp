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

func (ms *MessageStore) Create(msg *model.Message) error {
	sqlStmt := `
	INSERT INTO messages (clientId, username, content, time) VALUES (?, ?, ?, ?)
	`
	_, err := ms.DB.Exec(sqlStmt, msg.ClientID, msg.Username, msg.Data, msg.Time.UnixMilli())
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (ms *MessageStore) GetLast(count int) ([]*model.Message, error) {
	messages := []*model.Message{}

	rows, err := ms.DB.Query("SELECT * FROM messages ORDER BY time ASC LIMIT ?", count)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return messages, err
	}

	for rows.Next() {
		var clientId string
		var username string
		var content string
		var msec int64

		err = rows.Scan(&clientId, &username, &content, &msec)
		if err != nil {
			log.Println(err)
			return messages, err
		}

		messages = append(messages, &model.Message{
			ClientID: clientId,
			Username: username,
			Data:     content,
			Time:     time.UnixMilli(msec),
		})
	}

	return messages, nil
}
