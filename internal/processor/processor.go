package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func FindFunctionSignatures(bytecode []byte) []string {

	interpreter := BytecodeInterpreter{
		Pc:         0,
		Bytecode:   bytecode,
		Signatures: []string{},
		JumpDests:  [][2]byte{},
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

	//We remove 1 to the program counter not to lose the DUP1 that was consumed above
	if err := interpreter.AdjustPc(interpreter.Pc - 1); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		nextBytes, err := interpreter.NextBytes(DISPATCHER_SEQUENCE_LENGTH)

		if err == io.EOF {
			break
		}

		//We will check on DUP1,PUSH4 and PUSH2 to confirm that we are indeed in a dispatcher sequence
		dup1 := nextBytes[0]

		if dup1 != byte(vm.DUP1) {
			break
		}

		push4 := nextBytes[1]

		if push4 != byte(vm.PUSH4) {
			break
		}

		sig := nextBytes[2:6]

		eqOrGt := nextBytes[6]

		push2 := nextBytes[7]

		if push2 != byte(vm.PUSH2) {
			break
		}

		var destination = [2]byte{nextBytes[8], nextBytes[9]}

		switch eqOrGt {
		case byte(vm.EQ):
			interpreter.Signatures = append(interpreter.Signatures, common.Bytes2Hex(sig))
		case byte(vm.GT):
			interpreter.JumpDests = append(interpreter.JumpDests, destination)
		default:
			//Went out of dispatcher sequence
			break
		}

	}

	//fmt.Println(interpreter.Signatures)
	return interpreter.Signatures
}
