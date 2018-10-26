package g

import (
	"flag"
	"os"

	dlog "github.com/open-falcon/falcon-plus/logger"
)

var (
	strategyCfg  = flag.String("s", "cfg/strategy.dev.json", "specify strategy json file")
	StrategyFile string
)

func InitStrategyFile() {
	flag.Parse()
	cfgFile := *strategyCfg
	if cfgFile == "" {
		dlog.Infof("strategy file not specified: use -c $filename")
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		dlog.Infof("strategy file specified not found:%s\n", cfgFile)
	}

	StrategyFile = cfgFile
	dlog.Infof("use strategy file : %s", StrategyFile)
}
