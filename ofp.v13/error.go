package ofp

import (
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

const (
	ET_HELLO_FAILED ErrorType = iota
	ET_BAD_REQUEST
	ET_BAD_ACTION
	ET_BAD_INSTRUCTION
	ET_BAD_MATCH
	ET_FLOW_MOD_FAILED
	ET_GROUP_MOD_FAILED
	ET_PORT_MOD_FAILED
	ET_TABLE_MOD_FAILED
	ET_QUEUE_OP_FAILED
	ET_SWITCH_CONFIG_FAILED
	ET_ROLE_REQUEST_FAILED
	ET_METER_MOD_FAILED
	ET_TABLE_FEATURES_FAILED
	ET_EXPERIMENTER ErrorType = 0xffff
)

type ErrorType uint16

const (
	HFC_INCOMPATIBLE ErrorCode = iota
	HFC_EPERM
)

const (
	BRC_BAD_VERSION ErrorCode = iota
	BRC_BAD_TYPE
	BRC_BAD_MULTIPART
	BRC_BAD_EXPERIMENTER
	BRC_BAD_EXP_TYPE
	BRC_EPERM
	BRC_BAD_LEN
	BRC_BUFFER_EMPTY
	BRC_BUFFER_UNKNOWN
	BRC_BAD_TABLE_ID
	BRC_IS_SLAVE
	BRC_BAD_PORT
	BRC_BAD_PACKET
	BRC_MULTIPART_BUFFER_OVERFLOW
)

const (
	BAC_BAD_TYPE ErrorCode = iota
	BAC_BAD_LEN
	BAC_BAD_EXPERIMENTER
	BAC_BAD_EXP_TYPE
	BAC_BAD_OUT_PORT
	BAC_BAD_ARGUMENT
	BAC_EPERM
	BAC_TOO_MANY
	BAC_BAD_QUEUE
	BAC_BAD_OUT_GROUP
	BAC_MATCH_INCONSISTENT
	BAC_UNSUPPORTED_ORDER
	BAC_BAD_TAG
	BAC_BAD_SET_TYPE
	BAC_BAD_SET_LEN
	BAC_BAD_SET_ARGUMENT
)

const (
	BIC_UNKNOWN_INST ErrorCode = iota
	BIC_UNSUP_INST
	BIC_BAD_TABLE_ID
	BIC_UNSUP_METADATA
	BIC_UNSUP_METADATA_MASK
	BIC_BAD_EXPERIMENTER
	BIC_BAD_EXP_TYPE
	BIC_BAD_LEN
	BIC_EPERM
)

const (
	// Unsupported match type specified by the match
	BMC_BAD_TYPE ErrorCode = iota

	// Length problem in match
	BMC_BAD_LEN

	// Match uses an unsupported tag/encap
	BMC_BAD_TAG

	// Unsupported datalink addr mask - switch does not support
	// arbitrary datalink address mask
	BMC_BAD_DL_ADDR_MASK

	// Unsupported network addr mask - switch does not support
	// arbitrary network address mask
	BMC_BAD_NW_ADDR_MASK

	// Unsupported combination of fields masked or omitted in the match.
	BMC_BAD_WILDCARDS

	// Unsupported field type in the match
	BMC_BAD_FIELD

	// Unsupported value in a match field
	BMC_BAD_VALUE

	// Unsupported mask specified in the match
	BMC_BAD_MASK

	// A prerequisite was not met
	BMC_BAD_PREREQ

	// A field type was duplicated
	BMC_DUP_FIELD

	// Permissions error
	BMC_EPERM
)

const (
	FMFC_UNKNOWN ErrorCode = iota
	FMFC_TABLE_FULL
	FMFC_BAD_TABLE_ID
	FMFC_OVERLAP
	FMFC_EPERM
	FMFC_BAD_TIMEOUT
	FMFC_BAD_COMMAND
	FMFC_BAD_FLAGS
)

const (
	GMFC_GROUP_EXISTS ErrorCode = iota
	GMFC_INVALID_GROUP
	GMFC_WEIGHT_UNSUPPORTED
	GMFC_OUT_OF_GROUPS
	GMFC_OUT_OF_BUCKETS
	GMFC_CHAINING_UNSUPPORTED
	GMFC_WATCH_UNSUPPORTED
	GMFC_LOOP
	GMFC_UNKNOWN_GROUP
	GMFC_CHAINED_GROUP
	GMFC_BAD_TYPE
	GMFC_BAD_COMMAND
	GMFC_BAD_BUCKET
	GMFC_BAD_WATCH
	GMFC_EPERM
)

const (
	PMFC_BAD_PORT ErrorCode = iota
	PMFC_BAD_HW_ADDR
	PMFC_BAD_CONFIG
	PMFC_BAD_ADVERTISE
	PMFC_EPERM
)

const (
	TMFC_BAD_TABLE ErrorCode = iota
	TMFC_BAD_CONFIG
	TMFC_EPERM
)

const (
	QOFC_BAD_PORT ErrorCode = iota
	QOFC_BAD_QUEUE
	QOFC_EPERM
)

const (
	SCFC_BAD_FLAGS ErrorCode = iota
	SCFC_BAD_LEN
	SCFC_EPERM
)

const (
	RRFC_STALE ErrorCode = iota
	RRFC_UNSUP
	RRFC_BAD_ROLE
)

const (
	MMFC_UNKNOWN ErrorCode = iota
	MMFC_METER_EXISTS
	MMFC_INVALID_METER
	MMFC_UNKNOWN_METER
	MMFC_BAD_COMMAND
	MMFC_BAD_FLAGS
	MMFC_BAD_RATE
	MMFC_BAD_BURST
	MMFC_BAD_BAND
	MMFC_BAD_BAND_VALUE
	MMFC_OUT_OF_METERS
	MMFC_OUT_OF_BANDS
)

const (
	TFFC_BAD_TABLE ErrorCode = iota
	TFFC_BAD_METADATA
	TFFC_BAD_TYPE
	TFFC_BAD_LEN
	TFFC_BAD_ARGUMENT
	TFFC_EPERM
)

type ErrorCode uint16

// There are times that the switch needs to notify the controller
// of a problem. This is done with the T_ERROR_MSG message
type ErrorMsg struct {
	// Type value indicates the high-level type of error.
	Type ErrorType

	// Code value is interpreted based on the type.
	Code ErrorCode

	// Data is variable length and interpreted based
	// on the type and code. Unless specified otherwise,
	// the data field contains at least 64 bytes of the
	// failed request that caused the error message to
	// be generated, if the failed request is shorter
	// than 64 bytes it should be the full request without any padding.
	Data []byte
}

func (e *ErrorMsg) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = encoding.ReadFrom(r, &e.Type, &e.Code)
	if err != nil {
		return
	}

	e.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return n + int64(len(e.Data)), nil
}

func (e *ErrorMsg) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e)
}

type ErrorExperimenterMsg struct {
	Type         ErrorType
	ExpType      uint16
	Experimenter uint32
	Data         []byte
}
