package ofp

type (
	pad  []byte
	pad1 [1]uint8
	pad2 [2]uint8
	pad3 [3]uint8
	pad4 [4]uint8
	pad5 [5]uint8
	pad6 [6]uint8
	pad7 [7]uint8
	pad8 [8]uint8
)

var (
	defaultPad  pad
	defaultPad1 pad1
	defaultPad2 pad2
	defaultPad3 pad3
	defaultPad4 pad4
	defaultPad5 pad5
	defaultPad6 pad6
	defaultPad7 pad7
	defaultPad8 pad8
)

// padLen returns a size of the padding for the given length.
func padLen(length int) int {
	return (length+7)/8*8 - length
}

// makePad creates a new padding based on the given length according
// to the formula: (length + 7) / 8 * 8 - length
func makePad(length int) []byte {
	return make([]byte, padLen(length))
}
