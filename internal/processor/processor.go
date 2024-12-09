package processor

import (
	"bytes"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

var dispatcherSequence = [5]byte{vm.DUP1, byte(vm.PUSH4), byte(vm.EQ), byte(vm.PUSH2), byte(vm.JUMPI)}

var dispatcherOpcodes = map[byte]bool{
	vm.DUP1:        true,
	byte(vm.PUSH4): true,
	byte(vm.EQ):    true,
	byte(vm.PUSH2): true,
	byte(vm.JUMPI): true,
}

func FindFunctionSignatures(bytecode []byte) []string {

	foundSignatures := *new([]string)

	reader := bytes.NewReader(bytecode)

	currentSequence := make([]byte, 0, 5)
	lastSignature := make([]byte, 4)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if len(currentSequence) == 5 {
			//If currentSequence == 5 then we know the last PUSH4 pushed a function signature
			foundSignatures = append(foundSignatures, common.Bytes2Hex(lastSignature))

			lastSignature = make([]byte, 4, 4)
			currentSequence = nil

		}

		if _, ok := dispatcherOpcodes[b]; !ok {

			if len(foundSignatures) > 0 {
				//exited function dispatcher
				break
			}
			continue
		}

		switch b {

		case dispatcherSequence[0], dispatcherSequence[1],
			dispatcherSequence[2], dispatcherSequence[3],
			dispatcherSequence[4]:

			currentSequence = append(currentSequence, b)

			switch b {

			case dispatcherSequence[1]:
				reader.Read(lastSignature)

			case dispatcherSequence[3]:
				//consuming the destination of JUMPI
				for i := 0; i < 2; i++ {
					reader.ReadByte()
				}

			}

		}

	}

	//fmt.Println(foundSignatures)
	return foundSignatures
}
