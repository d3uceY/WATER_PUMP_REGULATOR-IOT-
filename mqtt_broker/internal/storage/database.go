package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	// i mean, it is already a single writer but
	// i feel special when i just write unecessary shit
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	fmt.Println("DB initialized on Bro")

	return DB.Ping()
}

// creates the tables on server load
func CreateTables() {
	query := `
	CREATE TABLE IF NOT EXISTS telegram_chats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		chat_id INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_telegram_chats_chat_id
		 ON telegram_chats(chat_id);
     `

	_, err := DB.Exec(query)
	if err != nil {
		fmt.Printf("SQL Error: %v\nQuery: %s\n", err, query)
		panic(err)
	}
}

// runs migrations (obviously)
func RunMigrations() {
	// maybe when i get something to migrate lol (because most times i always need to)
}
