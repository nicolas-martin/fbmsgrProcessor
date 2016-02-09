package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type incomingMessage struct {
	ID         int
	AuthorID   string
	ThreadID   string
	AuthorName string
	Message    string
	isRead     bool
}

func main() {
	db, err := sql.Open("mysql", "root:password@/mydb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	message := getUnReadMessages(db)

	fmt.Printf("== %s ", message.AuthorName)

	//TODO: Process the new message
	var returnMessage string
	if message.Message == "hello" {
		returnMessage = "Good day sir."
	}

	toggleIsRead(db, message)

	if len(returnMessage) != 0 {
		insertMessageToSend(db, message, returnMessage)
	}

}

func getUnReadMessages(db *sql.DB) incomingMessage {
	stmtOut, err := db.Prepare("SELECT id, author_id, thread_id, author_name, message, isRead FROM inbox WHERE isRead = ?")
	if err != nil {
		panic(err.Error())
	}

	defer stmtOut.Close()

	var message incomingMessage

	//TODO: loop for multiple rows
	err = stmtOut.QueryRow(0).Scan(&message.ID, &message.AuthorID, &message.ThreadID, &message.AuthorName, &message.Message, &message.isRead)
	if err != nil {
		panic(err.Error())
	}

	return message
}

func toggleIsRead(db *sql.DB, message incomingMessage) {
	stmUpdateInbox, err := db.Prepare("update inbox set isRead = 1 where id = ?;")
	if err != nil {
		panic(err.Error())
	}

	defer stmUpdateInbox.Close()
	_, err = stmUpdateInbox.Exec(message.ID)
	if err != nil {
		panic(err.Error())
	}
}

func insertMessageToSend(db *sql.DB, message incomingMessage, returnMessage string) {
	stmNewOutbox, err := db.Prepare("INSERT INTO `outbox` (`author_id`, `thread_id`, `author_name`, `message`) values (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	defer stmNewOutbox.Close()
	fmt.Printf("returnMessage = %s", returnMessage)
	if err != nil {
		panic(err.Error())
	}
}
