package serve

import (
	"context"
	"net/http"

	"github.com/likecoin/likechain/app"
	"github.com/likecoin/likecoin-chain-tx-indexer/db"
	"github.com/likecoin/likecoin-chain-tx-indexer/logger"
	"github.com/likecoin/likecoin-chain-tx-indexer/poller"
	"github.com/likecoin/likecoin-chain-tx-indexer/rest"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "serve",
	Short: "Run the indexing service and expose HTTP API",
	Run: func(cmd *cobra.Command, args []string) {
		restConn, err := db.NewConnFromCmdArgs(cmd)
		if err != nil {
			logger.L.Panicw("Cannot connect to Postgres", "error", err)
		}
		defer restConn.Close(context.Background())
		pollerConn, err := db.NewConnFromCmdArgs(cmd)
		if err != nil {
			logger.L.Panicw("Cannot connect to Postgres", "error", err)
		}
		defer pollerConn.Close(context.Background())
		err = db.InitDB(pollerConn)
		if err != nil {
			logger.L.Panicw("Cannot initialize Postgres database", "error", err)
		}
		listenAddr, err := cmd.PersistentFlags().GetString("listen-addr")
		if err != nil {
			logger.L.Panicw("Cannot get listen address from command line parameters", "error", err)
		}
		lcdEndpoint, err := cmd.PersistentFlags().GetString("lcd-endpoint")
		if err != nil {
			logger.L.Panicw("Cannot get lcd endpoint address from command line parameters", "error", err)
		}
		ignoreHeightDiff, err := cmd.PersistentFlags().GetBool("ignore-height-difference")
		if err != nil {
			logger.L.Panicw("Cannot get ignore-height-difference param from command line parameters", "error", err)
		}
		if lcdEndpoint[len(lcdEndpoint)-1] == '/' {
			lcdEndpoint = lcdEndpoint[:len(lcdEndpoint)-1]
		}
		ctx := poller.CosmosCallContext{
			Codec:       app.MakeCodec(),
			Client:      &http.Client{},
			LcdEndpoint: lcdEndpoint,
		}

		if !ignoreHeightDiff {
			const heightDiffLimit = 10000
			dbHeight, err := db.GetLatestHeight(restConn)
			if err != nil {
				logger.L.Panicw("Cannot get height from database", "error", err)
			}
			blockResult, err := poller.GetBlock(&ctx, 0)
			lcdHeight := blockResult.Block.Height
			if lcdHeight-dbHeight > heightDiffLimit {
				logger.L.Fatalw("height difference larger than limit, please run `import` or add --ignore-height-difference flag", "db_height", dbHeight, "lcd_height", lcdHeight, "limit", heightDiffLimit)
			}
		}
		go rest.Run(restConn, listenAddr, lcdEndpoint)
		poller.Run(pollerConn, &ctx)
	},
}

func init() {
	Command.PersistentFlags().String("lcd-endpoint", "http://localhost:1317", "LikeCoin chain lite client RPC endpoint")
	Command.PersistentFlags().String("listen-addr", "localhost:8997", "HTTP API serving address")
	Command.PersistentFlags().Bool("ignore-height-difference", false, "start serving and polling without import even if the height lags behind too much")
}
