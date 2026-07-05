package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func getAppDataDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, "d3uc3y", "water_pump_regulator", "database", "store")
	os.MkdirAll(path, 0755)
	return path
}

func InitDB() error {

	path := getAppDataDir()
	fmt.Printf(path)
	DB, err := sql.Open("sqlite", path)

	if err != nil {
		return err
	}

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
