package ofp

const (
	DESC_STR_LEN   = 256
	SERIAL_NUM_LEN = 32
)

type Desc struct {
	Manufacturer string
	Hardware     string
	Software     string
	SerialNum    string
	Datapath     string
}
