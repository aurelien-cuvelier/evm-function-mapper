package processor

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Length in bytes of a dispatcher sequence from DUP1 to JUMPI
// The length is the same for both types
const DISPATCHER_SEQUENCE_LENGTH = 11

type BytecodeInterpreter struct {
	Bytecode   []byte
	Pc         int
	Signatures []string
	JumpDests  []int
}

func (bi *BytecodeInterpreter) ProcessNode() {

	for {

		nextBytes, err := bi.NextBytes(DISPATCHER_SEQUENCE_LENGTH)

		if err != nil {

			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				os.Exit(1)
			}

		}

		//We will check on DUP1,PUSH4 and PUSH2 to exit when we get out of a sequencer pattern and break
		dup1 := nextBytes[0]
		if dup1 != byte(vm.DUP1) {
			break
		}

		push4 := nextBytes[1]
		if push4 != byte(vm.PUSH4) {
			break
		}

		eqOrGt := nextBytes[6]

		push2 := nextBytes[7]
		if push2 != byte(vm.PUSH2) {
			break
		}

		switch eqOrGt {
		case byte(vm.GT):
			{
				//We are converting the 2 bytes of the next JUMPDEST pc to int
				//for convenience indexing the bytecode
				binDest := [2]byte{nextBytes[8], nextBytes[9]}
				var destination int = int(binDest[0])<<8 | int(binDest[1])

				bi.JumpDests = append(bi.JumpDests, destination)
			}
		case byte(vm.EQ):
			{
				//We only push the signature in case we are in a linear sequence
				//otherwise it will be duplicated later
				sig := nextBytes[2:6]
				bi.Signatures = append(bi.Signatures, common.Bytes2Hex(sig))
			}
		default:
			//Should not be possible to come here, but for safety purposes
			fmt.Printf("Unexpected bytecode, got 0x%x but expected 0x%x or 0x%x\n", eqOrGt, byte(vm.GT), byte(vm.EQ))
			break
		}

	}

	return
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
			//Just consuming the pushed bytes
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
