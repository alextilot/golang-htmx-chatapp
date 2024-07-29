package store

import (
	"database/sql"
	"log"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/model"
)

type MessageStore struct {
	DB *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
	store := &MessageStore{DB: db}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS messages (username text not null, content text, time INTEGER);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return store
}

func (ms *MessageStore) Create(msg *model.Message) (*model.Message, error) {
	sqlStmt := `
	INSERT INTO messages (username, content, time) VALUES (?, ?, ?)
	`
	sqlResult, err := ms.DB.Exec(sqlStmt, msg.Username, msg.Data, msg.Time.UnixMilli())
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	msg.Number = id
	return msg, nil
}

func (ms *MessageStore) GetRows(count int, offset int) ([]*model.Message, error) {
	messages := []*model.Message{}

	rows, err := ms.DB.Query("SELECT row_number() over() AS ID, * FROM messages WHERE rowid < ? ORDER BY rowid DESC LIMIT ?", offset, count)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return messages, err
	}

	for rows.Next() {
		var rowid int64
		var username string
		var content string
		var msec int64

		err = rows.Scan(&rowid, &username, &content, &msec)
		if err != nil {
			log.Println(err)
			return messages, err
		}

		messages = append(messages, &model.Message{
			Number:   rowid,
			Username: username,
			Data:     content,
			Time:     time.UnixMilli(msec),
		})
	}

	return messages, nil
}

func (ms *MessageStore) GetMostRecent(count int) ([]*model.Message, error) {
	messages := []*model.Message{}

	rows, err := ms.DB.Query("SELECT row_number() over() AS ID, * FROM messages ORDER BY rowid DESC LIMIT ?", count)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return messages, err
	}

	for rows.Next() {
		var rowid int64
		var username string
		var content string
		var msec int64

		err = rows.Scan(&rowid, &username, &content, &msec)
		if err != nil {
			log.Println(err)
			return messages, err
		}

		messages = append(messages, &model.Message{
			Number:   rowid,
			Username: username,
			Data:     content,
			Time:     time.UnixMilli(msec),
		})
	}

	return messages, nil
}
