package processor

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/core/vm"
)

// Length in bytes of a dispatcher sequence from DUP1 to JUMPI
// The length is the same for both types
const DISPATCHER_SEQUENCE_LENGTH = 11

type BytecodeInterpreter struct {
	Bytecode   []byte
	Pc         int
	Signatures []string
	JumpDests  [][2]byte
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

func (bi *BytecodeInterpreter) AdjustPc(newPos int) error {

	if newPos < 0 || newPos > len(bi.Bytecode)-1 {
		return errors.New("New position outside of boundaries")
	}

	bi.Pc = newPos

	return nil
}
