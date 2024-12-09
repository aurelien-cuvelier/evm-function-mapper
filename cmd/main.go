package main

import (
	"github.com/aurelien-cuvelier/evm-function-mapper/internal/config"
	"github.com/aurelien-cuvelier/evm-function-mapper/internal/fetcher"
)

func main() {

	conf := config.GetConfig()

	if len(conf.Bytecode) == 0 {
		fetcher.FetchBytecode(conf)
	}

}
