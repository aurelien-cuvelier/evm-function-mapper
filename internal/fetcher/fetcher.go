package fetcher

import (
	"context"
	"fmt"
	"os"

	"github.com/aurelien-cuvelier/evm-function-mapper/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FetchBytecode(conf *config.Config) {

	client, err := ethclient.DialContext(context.Background(), conf.Rpc)

	if err != nil {
		fmt.Println("Could not instantiate ethereum client: ", err)
		os.Exit(1)
	}

	chainId, err := client.ChainID(context.Background())

	if err != nil {
		fmt.Println("Could not fetch chain id, is your RPC alive?")
		os.Exit(1)
	}

	bytecode, err := client.CodeAt(context.Background(), conf.Address, nil)

	if err != nil {
		fmt.Println("Could not fetch code: ", err)
		os.Exit(1)
	}

	fmt.Printf("Fetched %d bytes of code at %s on chain id %d\n", len(bytecode), conf.Address, chainId)

	conf.Bytecode = bytecode

}
