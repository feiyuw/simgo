package main

import (
	"flag"
	"simgo/logger"
	"simgo/ops"
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
