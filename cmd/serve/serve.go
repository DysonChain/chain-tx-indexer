package serve

import (
	"github.com/spf13/cobra"

	"github.com/likecoin/likecoin-chain-tx-indexer/pubsub"
	"github.com/likecoin/likecoin-chain-tx-indexer/rest"
)

var Command = &cobra.Command{
	Use:   "serve",
	Short: "Run the indexing service and expose HTTP API",
	Long:  "Deprecated. Use the `rest` and `poll` subcommands to run HTTP API server and poller separately instead.",
	Run: func(cmd *cobra.Command, args []string) {
		go ServeHTTP(cmd)
		ServePoller(cmd)
	},
}

func init() {
	Command.AddCommand(PollerCommand, HTTPCommand)
	rest.ConfigCmd(Command)
	pubsub.ConfigCmd(Command)
}
