package utility

import (
	"aria/backend/database"
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"sync"
)

var dbMu sync.Mutex

func GetDB() *sql.DB {
	dbMu.Lock()
	defer dbMu.Unlock()

	home, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fp := home + "/app.db"
	db, err := sql.Open("sqlite", fp)

	if err != nil {
		panic(err)
	}

	return db
}

var DatabaseMiddleware = NewMiddleware(
	func(c *Context) {
		db := GetDB()
		c.Queries = database.New(db)
	},
)
