package main

import (
	"os"

	"github.com/pandakn/cafe-beans/config"
	"github.com/pandakn/cafe-beans/modules/servers"
	"github.com/pandakn/cafe-beans/pkg/database"
)

func envPath() string {
	// e.g., ./command-line-arguments a b c d (5 arguments)
	// os.Args[1] is a

	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DbConnect(cfg.Db())

	servers.NewServer(cfg, db).Start()

	defer db.Close()
}
