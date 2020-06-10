package importdb

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/likecoin/likecoin-chain-tx-indexer/db"
	"github.com/likecoin/likecoin-chain-tx-indexer/importdb"
	"github.com/likecoin/likecoin-chain-tx-indexer/logger"
)

var Command = &cobra.Command{
	Use:   "import",
	Short: "Import from existing LikeCoin chain database",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := db.NewConnFromCmdArgs(cmd)
		if err != nil {
			logger.L.Panicw("Cannot connect to Postgres", "error", err)
		}
		defer conn.Close(context.Background())
		err = db.InitDB(conn)
		if err != nil {
			logger.L.Panicw("Cannot initialize Postgres database", "error", err)
		}
		likedPath, err := cmd.PersistentFlags().GetString("liked-path")
		if err != nil {
			logger.L.Panicw("Cannot get liked data folder path from command line parameters", "error", err)
		}
		importdb.Run(conn, likedPath)
	},
}

func Execute() {
	if err := Command.Execute(); err != nil {
		logger.L.Fatalw("Import command execution failed", "error", err)
	}
}

func init() {
	Command.PersistentFlags().String("liked-path", "./.liked", "location of the LikeCoin chain database folder (.liked)")
}
