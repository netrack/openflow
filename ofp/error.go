package ofp

import (
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/encoding"
)

const (
	ErrorTypeHelloFailed ErrorType = iota
	ErrorTypeBadRequest
	ErrorTypeBadAction
	ErrorTypeBadInstruction
	ErrorTypeBadMatch
	ErrorTypeFlowModFailed
	ErrorTypeGroupModFailed
	ErrorTypePortModFailed
	ErrorTypeTableModFailed
	ErrorTypeQueueOpFailed
	ErrorTypeSwitchConfigFailed
	ErrorTypeRoleRequestFailed
	ErrorTypeMeterModFailed
	ErrorTypeTableFeaturesFailed
	ErrorTypeExperimenter ErrorType = 0xffff
)

type ErrorType uint16

const (
	ErrorCodeHelloFailedIncompatible ErrorCode = iota
	ErrorCodeHelloFailedPermissions
)

type ErrorCode uint16

const (
	// Unsupported match type specified by the match
	ErrorCodeBadMatchBadType ErrorCode = iota

	// Length problem in match
	ErrorCodeBadMatchBadLen

	// Match uses an unsupported tag/encap
	ErrorCodeBadMatchBadTag

	// Unsupported datalink addr mask - switch does not support
	// arbitrary datalink address mask
	ErrorCodeBadMatchBadDatalinkAddressMask

	// Unsupported network addr mask - switch does not support
	// arbitrary network address mask
	ErrorCodeBadMatchBadNetworkAddressMask

	// Unsupported combination of fields masked or omitted in the match.
	ErrorCodeBadMatchBadWildcards

	// Unsupported field type in the match
	ErrorCodeBadMatchBadField

	// Unsupported value in a match field
	ErrorCodeBadMatchBadValue

	// Unsupported mask specified in the match
	ErrorCodeBadMatchBadMask

	// A prerequisite was not met
	ErrorCodeBadMatchBadPrereq

	// A field type was duplicated
	ErrorCodeBadMatchDuplicateField

	// Permissions error
	ErrorCodeBadMatchPermissions
)

const (
	// ErrorCodeBadRequestBadVersion ...
	ErrorCodeBadRequestBadVersion ErrorCode = iota
	ErrorCodeBadRequestBadType
	ErrorCodeBadRequestBadMultipart
	ErrorCodeBadRequestBadExperimenter
	ErrorCodeBadRequestBadExpType
	ErrorCodeBadRequestPermissions
	ErrorCodeBadRequestBadLen
	ErrorCodeBadRequestBufferEmpty
	ErrorCodeBadRequestBufferUnknown
	ErrorCodeBadRequestBadTableID
	ErrorCodeBadRequestIsSlave
	ErrorCodeBadRequestBadPort
	ErrorCodeBadRequestBadPacket
	ErrorCodeBadRequestMultipartBufferOverflow
)

const (
	ErrorCodeBadInstructionUnknownInstruction ErrorCode = iota
	ErrorCodeBadInstructionUnsupportedInstruction
	ErrorCodeBadInstructionBadTableID
	ErrorCodeBadInstructionUnsupportedMetadata
	ErrorCodeBadInstructionUnsupportedMetadataMask
	ErrorCodeBadInstructionBadExperimenter
	ErrorCodeBadInstructionBadExpType
	ErrorCodeBadInstructionBadLen
	ErrorCodeBadInstructionPermissions
)

const (
	ErrorCodeFlowModFailedUnknown ErrorCode = iota
	ErrorCodeFlowModFailedTableFull
	ErrorCodeFlowModFailedBadTableID
	ErrorCodeFlowModFailedOverlap
	ErrorCodeFlowModFailedPermissions
	ErrorCodeFlowModFailedBadTimeout
	ErrorCodeFlowModFailedBadCommand
	ErrorCodeFlowModFailedBadFlags
)

const (
	ErrorCodeGroupModFailedGroupExists ErrorCode = iota
	ErrorCodeGroupModFailedInvalidGrop
	ErrorCodeGroupModFailedWeightUnsupported
	ErrorCodeGroupModFailedOutOfGroups
	ErrorCodeGroupModFailedOutOfBuckets
	ErrorCodeGroupModFailedChainingUnsupported
	ErrroCodeGroupModFailedWatchUnsupported
	ErrorCodeGroupModFailedLoop
	ErrorCodeGroupModFailedUnknownGroup
	ErrorCodeGroupModFailedChainedGroup
	ErrorCodeGroupModBadType
	ErrorCodeGroupModBadCommand
	ErrorCodeGroupModBadBucket
	ErrorCodeGroupModBadWatch
	ErrorCodeGroupModPermissions
)

const (
	ErrorCodePortModFailedBadPort ErrorCode = iota
	ErrorCodePortModFailedBadHardwareAddress
	ErrorCodePortModFailedBadConfig
	ErrorCodePortModFailedBadAdvertise
	ErrorCodePortModFailedPermissions
)

const (
	ErrorCodeTableModFailedBadTable ErrorCode = iota
	ErrorCodeTableModFailedBadConfig
	ErrorCodeTableModFailedPermissions
)

const (
	ErrorCodeQueueOpFailedBadPort ErrorCode = iota
	ErrorCodeQueueOpFailedBadQueue
	ErrorCodeQueueOpFailedPermissions
)

const (
	ErrorCodeSwitchConfigFailedBadFlags ErrorCode = iota
	ErrorCodeSwitchConfigFailedBadLen
	ErrorCodeSwitchConfigFailedPermissions
)

const (
	ErrorCodeRoleRequestFailedStale ErrorCode = iota
	ErrorCodeRoleRequestFailedUnsup
	ErrorCodeRoleRequestFailedBadRole
)

const (
	ErrorCodeMeterModFailedUnknown ErrorCode = iota
	ErrorCodeMeterModFailedMeterExists
	ErrorCodeMeterModFailedInvalidMeter
	ErrorCodeMeterModFailedUnknownMeter
	ErrorCodeMeterModFailedBadCommand
	ErrorCodeMeterModFailedBadFlags
	ErrorCodeMeterModFailedBadRate
	ErrorCodeMeterModFailedBadBurst
	ErrorCodeMeterModFailedBadBand
	ErrorCodeMeterModFailedBadBandValue
	ErrorCodeMeterModFailedOutOfMeters
	ErrorCodeMeterModFailedOutOfBands
)

const (
	ErrorCodeTypeTableFeaturesFailedBadTable ErrorCode = iota
	ErrorCodeTypeTableFeaturesFailedBadMetadata
	ErrorCodeTypeTableFeaturesFailedBadType
	ErrorCodeTypeTableFeaturesFailedBadLen
	ErrorCodeTypeTableFeaturesFailedBadArgument
	ErrorCodeTypeTableFeaturesFailedPermissions
)

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

func (e *ErrorMsg) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e.Type, e.Code, e.Data)
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

type ErrorExperimenterMsg struct {
	Type         ErrorType
	ExpType      uint16
	Experimenter uint32
	Data         []byte
}

func (e *ErrorExperimenterMsg) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e.Type, e.ExpType, e.Experimenter, e.Data)
}

func (e *ErrorExperimenterMsg) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = encoding.ReadFrom(r, &e.Type, &e.ExpType, &e.Experimenter)
	if err != nil {
		return n, err
	}

	e.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return n + int64(len(e.Data)), nil
}
