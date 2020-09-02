package main

import (
	"flag"

	"github.com/feiyuw/simgo/logger"
	"github.com/feiyuw/simgo/ops"
)

var (
	addr = ""
)

func main() {
	flag.StringVar(&addr, "addr", ":1777", "OPS addr")
	flag.Parse()
	logger.Infof("main", "start OPS on %s", addr)
	ops.Start(addr)
}
