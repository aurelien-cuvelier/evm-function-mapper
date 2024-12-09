package config

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

type Config struct {
	Bytecode []byte
	Rpc      string
	Address  string
}

func GetConfig() *Config {

	var conf *Config
	var rootCmd = &cobra.Command{
		Use:   "evm-function-finder",
		Short: "A simple tool to try and find known function signature in EVM bytecode",
		Long:  "A simple tool to try and find known function signature in EVM bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedConf, err := parseFlags(cmd)
			if err != nil {
				return err
			}

			conf = parsedConf

			return nil
		},
	}

	rootCmd.Flags().String("rpc", "", "The RPC endpoint to fetch from")
	rootCmd.Flags().String("bytecode", "", "The compiled bytecode")
	rootCmd.Flags().String("address", "", "The contract address")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return conf

}

func parseFlags(cmd *cobra.Command) (*Config, error) {

	rpc, err := cmd.Flags().GetString("rpc")
	if err != nil {
		return nil, fmt.Errorf("error parsing rpc: %s\n", err)
	}
	if len(rpc) > 0 {

	}

	bytecode, err := cmd.Flags().GetString("bytecode")
	if err != nil {
		return nil, fmt.Errorf("error parsing bytecode: %s\n", bytecode)
	}

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		return nil, fmt.Errorf("error parsing address %s\n")
	}

	if bytecode == "" && (rpc == "" || address == "") {
		return nil, fmt.Errorf("you need to provide the bytecode if you don't provide the RPC & address\n")
	}

	return &Config{
		Rpc:      rpc,
		Bytecode: common.Hex2Bytes(bytecode),
		Address:  address,
	}, nil
}
