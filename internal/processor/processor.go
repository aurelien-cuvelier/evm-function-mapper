package processor

import (
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

var dispatcherOpcodes = map[byte]bool{
	vm.DUP1:        true,
	byte(vm.PUSH4): true,
	byte(vm.EQ):    true,
	byte(vm.PUSH2): true,
	byte(vm.JUMPI): true,
}

type BytecodeInterpreter struct {
	Bytecode  []byte
	Pc        int
	LastPush4 string
}

func (bi *BytecodeInterpreter) IsPushOpCode(op byte) (bool, int) {

	if op >= byte(vm.PUSH0) && op <= byte(vm.PUSH32) {
		//PUSH0 = 0x5f and is incremented by 1 for each additionnal pushed byte
		return true, int(op - 0x5f)
	}

	return false, 0
}

func (bi *BytecodeInterpreter) NextByte() (byte, error) {

	if len(bi.Bytecode) < bi.Pc-1 {
		return 0, io.EOF
	}
	nextByte := bi.Bytecode[bi.Pc]
	bi.Pc++

	return nextByte, nil
}

func (bi *BytecodeInterpreter) NextBytes(n int) ([]byte, error) {

	if bi.Pc+n > len(bi.Bytecode)-1 {
		return nil, io.EOF
	}

	nextBytes := bi.Bytecode[bi.Pc : bi.Pc+n]
	bi.Pc += n

	return nextBytes, nil
}

func (bi *BytecodeInterpreter) ReadUntil(op vm.OpCode) error {
	//Reads the bytecode until finding the first match of the given opcode
	for {

		nextByte, err := bi.NextByte()
		if err != nil {
			return err
		}
		isPush, size := bi.IsPushOpCode(nextByte)

		if nextByte == byte(op) {
			return nil
		}

		if isPush {
			_, err := bi.NextBytes(size)
			if err != nil {
				return err
			}
		}

	}

}

func FindFunctionSignatures(bytecode []byte) []string {

	interpreter := BytecodeInterpreter{
		Pc:        0,
		Bytecode:  bytecode,
		LastPush4: "",
	}

	foundSignatures := *new([]string)
	delimiters := [2]vm.OpCode{vm.CALLDATALOAD, vm.DUP1}

	for _, delimiter := range delimiters {
		//Consuming opcodes until we found the first DUP1 after CALLDATALOAD
		err := interpreter.ReadUntil(delimiter)

		if err != nil {
			if err == io.EOF {
				return foundSignatures
			} else {
				fmt.Println("Unknown error consuming delimiters: ", err)
				os.Exit(1)
			}

		}
	}

	for {

		nextByte, err := interpreter.NextByte()

		if err == io.EOF {
			break
		}

		if _, ok := dispatcherOpcodes[nextByte]; !ok {
			break
		}

		isPush, size := interpreter.IsPushOpCode(nextByte)

		if isPush {

			pushedData, err := interpreter.NextBytes(size)

			if err == io.EOF {
				break
			}

			switch size {
			case 4:

				nextByte, err = interpreter.NextByte()

				if nextByte == byte(vm.EQ) {
					interpreter.LastPush4 = common.Bytes2Hex(pushedData)
					foundSignatures = append(foundSignatures, interpreter.LastPush4)
				}

			}

		}

	}

	//fmt.Println(foundSignatures)
	return foundSignatures
}
