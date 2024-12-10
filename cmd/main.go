package main

import (
	"fmt"
	"os"

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

	_, err := fetcher.QueryEthSignatureDatabase(foundSignatures)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
