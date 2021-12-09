package daemon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/fetcher"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/zmq"
)

var DaemonCommand = cli.Command{
	Name:   "daemon",
	Action: runDaemon,
	Usage:  "run daemon",
}

func runDaemon(ctx *cli.Context) {

	exitSignals := make(chan os.Signal, 1)
	shouldExit := make(chan bool, 1)
	noError := true

	signal.Notify(exitSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go handleExitSig(exitSignals, shouldExit)

	redisPool, err := database.SetupRedis()
	if err != nil {
		logger.Error.Printf("Failed to setup Redis database connection: %s", err.Error())
		shouldExit <- true
		noError = false
	}

	if config.GetBool("zmq.saveMempoolEntries.enable") {
		sqlDB, err := database.SetupSQLite()
		if err != nil {
			logger.Error.Printf("Failed to setup SQLite database connection: %s", err.Error())
			shouldExit <- true
			noError = false
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				logger.Error.Println("sqlDB close err:", err.Error())
			}
		}()
	}

	if noError {
		startWorkers(redisPool)
	}

	<-shouldExit // wait till memod should exit
	logger.Info.Println("Memod exiting")
}

// handles exit signals
func handleExitSig(exitSignals chan os.Signal, shouldExit chan bool) {
	sig := <-exitSignals
	logger.Info.Println("Received signal", sig)
	shouldExit <- true
}

func startWorkers(pool *database.RedisPool) {

	if config.GetBool("mempool.enable") {
		logger.Info.Println("Starting with mempool fetching enabled")
		go fetcher.SetupMempoolFetcher(pool)
	}

	if config.GetBool("feeratefetcher.enable") {
		logger.Info.Println("Starting with feerate API fetching enabled")
		go fetcher.SetupFeerateAPIFetcher(pool)
	}

	if config.GetBool("zmq.enable") {
		logger.Info.Println("Starting with ZMQ interface enabled")
		go zmq.SetupZMQ(pool)
	}

}
