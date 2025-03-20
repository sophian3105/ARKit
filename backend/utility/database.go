package utility

import (
	"aria/backend/database"
	"database/sql"
	_ "embed"
	_ "modernc.org/sqlite"
	"os"
	"sync"
)

var dbMu sync.Mutex
var db *sql.DB

func GetDB() *sql.DB {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return db
	}

	home, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fp := home + "/app.db"
	db, err := sql.Open("sqlite", fp)

	if err != nil {
		panic(err)
	}

	// Add the schema to the database if it doesn't exist
	if _, err := db.Exec(database.DDL); err != nil {
		panic(err)
	}

	return db
}

var DatabaseMiddleware = NewMiddleware(
	func(c *Context) {
		c.Queries = database.New(GetDB())
	},
)
