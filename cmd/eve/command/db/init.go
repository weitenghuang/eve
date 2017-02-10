package db

import (
	"github.com/concur/rohr/service/rethinkdb"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	MAX_RETRY int = 60
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "To initialize eve db",
	Long:  `To initialize eve db`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Initializing DB...")
		db := connectDB(0)
		if err := db.Initialization(); err != nil {
			log.Panicln(err)
		}
		log.Println("DB is in good condition.")
	},
}

func connectDB(retry int) *rethinkdb.DbSession {
	db := rethinkdb.DefaultSession()
	if db == nil && retry < MAX_RETRY {
		retry++
		log.Println("Wait for", retry, "seconds to re-try rethinkdb server connection.")
		time.Sleep(time.Duration(retry) * time.Second)
		db = connectDB(retry)
	}
	return db
}
