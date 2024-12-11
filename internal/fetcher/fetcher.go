package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aurelien-cuvelier/evm-function-mapper/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FetchBytecode(conf *config.Config) {

	client, err := ethclient.DialContext(context.Background(), conf.Rpc)
	if err != nil {
		fmt.Println("Could not instantiate ethereum client: ", err)
		os.Exit(1)
	}

	defer client.Close()

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

type FourByteApiResponse struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		ID             int       `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		TextSignature  string    `json:"text_signature"`
		HexSignature   string    `json:"hex_signature"`
		BytesSignature string    `json:"bytes_signature"`
	} `json:"results"`
}

func QueryEthSignatureDatabase(signatures []string) ([]string, error) {

	for _, signature := range signatures {

		res, err := http.Get(fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/?hex_signature=%s", signature))

		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintln("Status code: ", res.StatusCode))
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var marshalledData FourByteApiResponse

		if err := json.Unmarshal(body, &marshalledData); err != nil {
			return nil, err
		}

		fmt.Printf("==================0x%s===================\n", signature)
		for _, entry := range marshalledData.Results {
			fmt.Printf("%s\n\n", entry.TextSignature)
		}

		res.Body.Close()

	}

	return []string{}, nil
}
