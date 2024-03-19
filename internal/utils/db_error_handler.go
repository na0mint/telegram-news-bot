package utils

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func HandleCloseDbConnection(conn *sqlx.Conn) {
	err := conn.Close()
	if err != nil {
		log.Printf("[ERROR] Failed to close connection to database: %v", err)
	}
}
