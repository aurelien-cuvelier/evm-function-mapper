package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/core/vm"
)

func FindFunctionSignatures(bytecode []byte) []string {

	interpreter := BytecodeInterpreter{
		Pc:         0,
		Bytecode:   bytecode,
		Signatures: []string{},
		JumpDests:  []int{},
	}

	delimiters := [2]vm.OpCode{vm.CALLDATALOAD, vm.DUP1}

	for _, delimiter := range delimiters {
		//Consuming opcodes until we found the first DUP1 after CALLDATALOAD
		if err := interpreter.ReadUntil(delimiter); err != nil {
			if err == io.EOF {
				return interpreter.Signatures
			} else {
				fmt.Println("Unknown error consuming delimiters: ", err)
				os.Exit(1)
			}
		}

	}

	//We simulate the pc of our 1st DUP1 -1 to be our 1st JUMPDEST
	interpreter.JumpDests = append(interpreter.JumpDests, interpreter.Pc-2)

	for {

		if len(interpreter.JumpDests) == 0 {
			//Iterating until we visited all the nodes JUMPDEST
			break
		}

		jumpDest := interpreter.JumpDests[0]

		interpreter.JumpDests = interpreter.JumpDests[1:]

		//We add one the the pc to skip the JUMPDEST byte and go straight to the following DUP1
		if err := interpreter.AdjustPc(jumpDest + 1); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		interpreter.ProcessNode()

	}

	//fmt.Println(interpreter.Signatures)
	return interpreter.Signatures
}
