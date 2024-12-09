package main

import (
	"github.com/aurelien-cuvelier/evm-function-mapper/internal/config"
	"github.com/aurelien-cuvelier/evm-function-mapper/internal/fetcher"
	"github.com/aurelien-cuvelier/evm-function-mapper/internal/processor"
)

func main() {

	conf := config.GetConfig()

	if len(conf.Bytecode) == 0 {
		fetcher.FetchBytecode(conf)
	}

	foundSignatures := processor.FindFunctionSignatures(conf.Bytecode)

	fetcher.QueryEthSignatureDatabase(foundSignatures)
}
