package ofp

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/netrack/openflow/internal/encoding"
)

// ErrType indicates high-level type of error.
type ErrType uint16

// ErrCode indicates the precise type of error. The value is
// interpreted based on the error type.
type ErrCode uint16

const (
	// ErrTypeHelloFailed is returned when hello protocol failed.
	ErrTypeHelloFailed ErrType = iota

	// ErrTypeBadRequest is returned when request was no understood.
	ErrTypeBadRequest

	// ErrTypeBadAction is returned when error is in action
	// description.
	ErrTypeBadAction

	// ErrTypeBadInstruction is returned when error is in
	// instruction list.
	ErrTypeBadInstruction

	// ErrTypeBadMatch is returned when error is in match.
	ErrTypeBadMatch

	// ErrTypeFlowModFailed is returned when problem is in modifying
	// flow entry.
	ErrTypeFlowModFailed

	// ErrTypeGroupModFailed is returned when problem is in modifying
	// group entry.
	ErrTypeGroupModFailed

	// ErrTypePortModFailed is returned when port modification
	// request failed.
	ErrTypePortModFailed

	// ErrTypeTableModFailed is returned when table modification
	// request failed.
	ErrTypeTableModFailed

	// ErrTypeQueueOpFailed is returned when queue operation failed.
	ErrTypeQueueOpFailed

	// ErrTypeSwitchConfigFailed is returned when switch configuration
	// request failed.
	ErrTypeSwitchConfigFailed

	// ErrTypeRoleRequestFailed is returned when controller role
	// request failed.
	ErrTypeRoleRequestFailed

	// ErrTypeMeterModFailed is returned when error is in meter.
	ErrTypeMeterModFailed

	// ErrTypeTableFeaturesFailed is returned when setting table
	// features failed.
	ErrTypeTableFeaturesFailed

	// ErrTypeExperimenter is an experimenter error type.
	ErrTypeExperimenter ErrType = 0xffff
)

func (t ErrType) String() string {
	text, ok := errTypeText[t]
	if !ok {
		return fmt.Sprintf("ErrType(%d)", t)
	}
	return text
}

var errTypeText = map[ErrType]string{
	ErrTypeHelloFailed:         "ErrTypeHelloFailed",
	ErrTypeBadRequest:          "ErrTypeBadRequest",
	ErrTypeBadAction:           "ErrTypeBadAction",
	ErrTypeBadInstruction:      "ErrTypeBadInstruction",
	ErrTypeBadMatch:            "ErrTypeBadMatch",
	ErrTypeFlowModFailed:       "ErrTypeFlowModFailed",
	ErrTypeGroupModFailed:      "ErrTypeGroupModFailed",
	ErrTypePortModFailed:       "ErrTypePortModFailed",
	ErrTypeTableModFailed:      "ErrTypeTableModFailed",
	ErrTypeQueueOpFailed:       "ErrTypeQueueOpFailed",
	ErrTypeSwitchConfigFailed:  "ErrTypeSwitchConfigFailed",
	ErrTypeRoleRequestFailed:   "ErrTypeRoleRequestFailed",
	ErrTypeMeterModFailed:      "ErrTypeMeterModFailed",
	ErrTypeTableFeaturesFailed: "ErrTypeTableFeaturesFailed",
	ErrTypeExperimenter:        "ErrTypeExperimenter",
}

const (
	// ErrCodeHelloFailedIncompatible is returned when there is no
	// compatible version, that switch or controller supports.
	ErrCodeHelloFailedIncompatible ErrCode = iota

	// ErrCodeHelloFailedPerm is returned when permission denied.
	ErrCodeHelloFailedPerm
)

const (
	// ErrCodeBadRequestBadVersion is returned when version in the message
	// header is not supported.
	ErrCodeBadRequestBadVersion ErrCode = iota

	// ErrCodeBadRequestBadType is returned when message type is
	// not supported.
	ErrCodeBadRequestBadType

	// ErrCodeBadRequestBadMultipart is returned when multipart request
	// type is not supported.
	ErrCodeBadRequestBadMultipart

	// ErrCodeBadRequestBadExperimenter is returned when experimenter
	// identifier is not supported.
	ErrCodeBadRequestBadExperimenter

	// ErrCodeBadRequestBadExpType is returned when experimenter message
	// type is not supported.
	ErrCodeBadRequestBadExpType

	// ErrCodeBadRequestPerm is returned when permission defined.
	ErrCodeBadRequestPerm

	// ErrCodeBadRequestLen is returned when wrong request length
	// for type was specified.
	ErrCodeBadRequestLen

	// ErrCodeBadRequestBufferEmpty is returned when specified buffer
	// has already been used.
	ErrCodeBadRequestBufferEmpty

	// ErrCodeBadRequestBufferUnknown is returned when specified buffer
	// does not exist.
	ErrCodeBadRequestBufferUnknown

	// ErrCodeBadRequestBadTableID is returned when specified table
	// identifier is invalid or does not exist.
	ErrCodeBadRequestBadTableID

	// ErrCodeBadRequestIsSlave is returned when permission denied
	// because controller is slave.
	ErrCodeBadRequestIsSlave

	// ErrCodeBadRequestBadPort is returned when invalid port specified.
	ErrCodeBadRequestBadPort

	// ErrCodeBadRequestBadPacket is returned when invalid port in
	// packet-out message specified.
	ErrCodeBadRequestBadPacket

	// ErrCodeBadRequestMultipartBufferOverflow is returned when
	// multipart request message overflowed the assigned buffer.
	ErrCodeBadRequestMultipartBufferOverflow
)

const (
	// ErrCodeBadActionType is returned when unknown action type was
	// specified.
	ErrCodeBadActionType ErrCode = iota

	// ErrCodeBadActionLen is returned when invalid length was specified.
	ErrCodeBadActionLen

	// ErrCodeBadActionExperimenter is returned when unknown experimenter
	// identifier was specified.
	ErrCodeBadActionExperimenter

	// ErrCodeBadActionExpType is returned when unknown action for
	// experimenter identifier was specified.
	ErrCodeBadActionExpType

	// ErrCodeBadActionOutPort is returned when the problem validating
	// out port has been encountered.
	ErrCodeBadActionOutPort

	// ErrCodeBadActionArgument is returned when invalid action argument
	// was specified.
	ErrCodeBadActionArgument

	// ErrCodeBadActionPerm is returned when permission denied.
	ErrCodeBadActionPerm

	// ErrCodeBadActionTooMany is returned when datapath cannot handle
	// this many actions.
	ErrCodeBadActionTooMany

	// ErrCodeBadActionQueue is returned when the problem validating
	// output queue have been encountered.
	ErrCodeBadActionQueue

	// ErrCodeBadActionOutGroup is returned when invalid group identifier
	// in forward action have been specified.
	ErrCodeBadActionOutGroup

	// ErrCodeBadActionMatchInconsistent is returned when action cannot
	// be applied for the specified match or set-field instruction is
	// missing prerequisite.
	ErrCodeBadActionMatchInconsistent

	// ErrCodeBadActionUnsupportedOrder is returned when action order is
	// unsupported for the action list in apply-actions instruction.
	ErrCodeBadActionUnsupportedOrder

	// ErrCodeBadActionTag is returned when actions use an unsupported
	// tag or encapsulation.
	ErrCodeBadActionTag

	// ErrCodeBadActionSetType is returned when unsupported type was
	// specified in set-field action.
	ErrCodeBadActionSetType

	// ErrCodeBadActionSetLen is returned when invalid length was
	// specified in set-field action.
	ErrCodeBadActionSetLen

	// ErrCodeBadActionSetArgument is returned bad argument was specified
	// in set-field action.
	ErrCodeBadActionSetArgument
)

const (
	// ErrCodeBadInstructionUnknown is returned when specified instruction
	// is unknown.
	ErrCodeBadInstructionUnknown ErrCode = iota

	// ErrCodeBadInstructionUnsupported is returned when switch or table
	// does not support the instruction.
	ErrCodeBadInstructionUnsupported

	// ErrCodeBadInstructionTableID is returned when invalid table
	// identifier was specified.
	ErrCodeBadInstructionTableID

	// ErrCodeBadInstructionUnsupportedMetadata is returned when
	// specified metadata unsupported by datapath.
	ErrCodeBadInstructionUnsupportedMetadata

	// ErrCodeBadInstructionUnsupportedMetadataMask is returned when
	// specified metadata mask unsupported by datapath.
	ErrCodeBadInstructionUnsupportedMetadataMask

	// ErrCodeBadInstructionExperimenter is returned when unknown
	// experimenter identifier was specified.
	ErrCodeBadInstructionExperimenter

	// ErrCodeBadInstructionExpType is returned when unknown instruction
	// for experimenter type was specified.
	ErrCodeBadInstructionExpType

	// ErrCodeBadInstructionLen is returned when wrong instruction
	// length was specified.
	ErrCodeBadInstructionLen

	// ErrCodeBadInstructionPerm is returned when permission denied.
	ErrCodeBadInstructionPerm
)

const (
	// ErrCodeBadMatchBadType is returned when unsupported match type
	// specified by the match.
	ErrCodeBadMatchBadType ErrCode = iota

	// ErrCodeBadMatchBadLen is returned when length problem in match.
	ErrCodeBadMatchBadLen

	// ErrCodeBadMatchBadTag is returned when match uses an unsupported
	// tag or encapsulation.
	ErrCodeBadMatchBadTag

	// ErrCodeBadMatchBadLinkMask is returned when unsupported datalink
	// address mask specified - switch does not support arbitrary
	// datalink address mask.
	ErrCodeBadMatchBadLinkMask

	// ErrCodeBadMatchBadNetMask is returned when unsupported network
	// address mask specified - switch does not support arbitrary
	// network address mask.
	ErrCodeBadMatchBadNetMask

	// ErrCodeBadMatchBadWildcards is returned when unsupported
	// combination of fields masked or omitted in the match.
	ErrCodeBadMatchBadWildcards

	// ErrCodeBadMatchBadField is returned when unsupported field type
	// in the match.
	ErrCodeBadMatchBadField

	// ErrCodeBadMatchBadValue is returned when unsupported value in a
	// match field.
	ErrCodeBadMatchBadValue

	// ErrCodeBadMatchBadMask is returned when unsupported mask
	// specified in the match.
	ErrCodeBadMatchBadMask

	// ErrCodeBadMatchBadPrereq is returned when a prerequisite was
	// not met.
	ErrCodeBadMatchBadPrereq

	// ErrCodeBadMatchDupField is returned when a field type was
	// duplicated.
	ErrCodeBadMatchDupField

	// ErrCodeBadMatchPerm is returned when permission denied.
	ErrCodeBadMatchPerm
)

const (
	// ErrCodeFlowModFailedUnknown is returned in case of unspecified error.
	ErrCodeFlowModFailedUnknown ErrCode = iota

	// ErrCodeFlowModFailedTableFull is returned when flow was not added
	// because table was full.
	ErrCodeFlowModFailedTableFull

	// ErrCodeFlowModFailedBadTableID is returned when table does not exist.
	ErrCodeFlowModFailedBadTableID

	// ErrCodeFlowModFailedOverlap is returned when it was attempted to add
	// overlapping flow with overlap checking flag set.
	ErrCodeFlowModFailedOverlap

	// ErrCodeFlowModFailedPerm is returned when permission denied.
	ErrCodeFlowModFailedPerm

	// ErrCodeFlowModFailedBadTimeout is returned when flow was not added
	// because of unsupported IDLE or hard timeout.
	ErrCodeFlowModFailedBadTimeout

	// ErrCodeFlowModFailedBadCommand is returned when unsupported or
	// unknown command was specified.
	ErrCodeFlowModFailedBadCommand

	// ErrCodeFlowModFailedBadFlags is returned when unsupported or unknown
	// flags where specified.
	ErrCodeFlowModFailedBadFlags
)

const (
	// ErrCodeGroupModFailedGroupExists is returned when group was not
	// added because a group addition operation attempted to replace an
	// already-present group.
	ErrCodeGroupModFailedGroupExists ErrCode = iota

	// ErrCodeGroupModFailedInvalidGroup is returned when group was not
	// added because invalid group identifier was specified.
	ErrCodeGroupModFailedInvalidGroup

	// ErrCodeGroupModFailedWeightUnsupported is returned when switch
	// does not support unequal load sharing with selected group.
	ErrCodeGroupModFailedWeightUnsupported

	// ErrCodeGroupModFailedOutOfGroups is returned when group table is
	// full.
	ErrCodeGroupModFailedOutOfGroups

	// ErrCodeGroupModFailedOutOfBuckets is returned when maximum number
	// of action buckets for a group has been exceed.
	ErrCodeGroupModFailedOutOfBuckets

	// ErrCodeGroupModFailedChainingUnsupported is returned when switch
	// does not support groups that forward to groups.
	ErrCodeGroupModFailedChainingUnsupported

	// ErrCodeGroupModFailedWatchUnsupported is returned when the specified
	// group cannot watch given port or group.
	ErrCodeGroupModFailedWatchUnsupported

	// ErrCodeGroupModFailedLoop is returned when group entry would cause
	// a loop.
	ErrCodeGroupModFailedLoop

	// ErrCodeGroupModFailedUnknownGroup is returned when group was not
	// modified because it does not exist.
	ErrCodeGroupModFailedUnknownGroup

	// ErrCodeGroupModFailedChainedGroup is return when group was not
	// deleted because another group is forwarding to it.
	ErrCodeGroupModFailedChainedGroup

	// ErrCodeGroupModBadType is returned when unsupported or unknown
	// group type was specified.
	ErrCodeGroupModBadType

	// ErrCodeGroupModBadCommand is returned when unsupported or unknown
	// command was specified.
	ErrCodeGroupModBadCommand

	// ErrCodeGroupModBadBucket is returned in case of error in bucket.
	ErrCodeGroupModBadBucket

	// ErrCodeGroupModBadWatch is returned in case of error in watch
	// port or group.
	ErrCodeGroupModBadWatch

	// ErrCodeGroupModPerm is returned when permission denied.
	ErrCodeGroupModPerm
)

const (
	// ErrCodePortModFailedBadPort is returned when specified port number
	// does not exist.
	ErrCodePortModFailedBadPort ErrCode = iota

	// ErrCodePortModFailedBadHwAddr is returned when specified hardware
	// address does not match the port number.
	ErrCodePortModFailedBadHwAddr

	// ErrCodePortModFailedBadConfig is returned when specified
	// configuration is invalid.
	ErrCodePortModFailedBadConfig

	// ErrCodePortModFailedBadAdvertise is returned when specified
	// advertise is invalid.
	ErrCodePortModFailedBadAdvertise

	// ErrCodePortModFailedPerm is returned when permission denied.
	ErrCodePortModFailedPerm
)

const (
	// ErrCodeTableModFailedBadTable is returned when specified table
	// does not exist.
	ErrCodeTableModFailedBadTable ErrCode = iota

	// ErrCodeTableModFailedBadConfig is returned when specified
	// configuration is invalid.
	ErrCodeTableModFailedBadConfig

	// ErrCodeTableModFailedPerm is returned when permission denied.
	ErrCodeTableModFailedPerm
)

const (
	// ErrCodeQueueOpFailedBadPort is returned when invalid port specified
	// or it does not exist.
	ErrCodeQueueOpFailedBadPort ErrCode = iota

	// ErrCodeQueueOpFailedBadQueue is returned when specified queue does
	// not exist.
	ErrCodeQueueOpFailedBadQueue

	// ErrCodeQueueOpFailedPerm is returned when permission denied.
	ErrCodeQueueOpFailedPerm
)

const (
	// ErrCodeSwitchConfigFailedBadFlags is returned when specified flags
	// are invalid.
	ErrCodeSwitchConfigFailedBadFlags ErrCode = iota

	// ErrCodeSwitchConfigFailedBadLen is returned when specified length
	// is invalid.
	ErrCodeSwitchConfigFailedBadLen

	// ErrCodeSwitchConfigFailedPerm is returned when permission denied.
	ErrCodeSwitchConfigFailedPerm
)

const (
	// ErrCodeRoleRequestFailedStale is returned when the message is stale.
	// Old generation identifier received.
	ErrCodeRoleRequestFailedStale ErrCode = iota

	// ErrCodeRoleRequestFailedUnsup is returned when controller role
	// change unsupported.
	ErrCodeRoleRequestFailedUnsup

	// ErrCodeRoleRequestFailedBadRole is returned when invalid role was
	// specified.
	ErrCodeRoleRequestFailedBadRole
)

const (
	// ErrCodeMeterModFailedUnknown is returned in case of unspecified error.
	ErrCodeMeterModFailedUnknown ErrCode = iota

	// ErrCodeMeterModFailedMeterExists is returned when meter not added
	// because it already exists.
	ErrCodeMeterModFailedMeterExists

	// ErrCodeMeterModFailedInvalidMeter is returned when specified meter
	// is invalid.
	ErrCodeMeterModFailedInvalidMeter

	// ErrCodeMeterModFailedUnknownMeter is returned when meter not
	// modified because it does not exist.
	ErrCodeMeterModFailedUnknownMeter

	// ErrCodeMeterModFailedBadCommand is returned when an unsupported
	// or unknown command was specified.
	ErrCodeMeterModFailedBadCommand

	// ErrCodeMeterModFailedBadFlags is returned when specified flag
	// configuration is unsupported.
	ErrCodeMeterModFailedBadFlags

	// ErrCodeMeterModFailedBadRate is returned when specified rate
	// is unsupported.
	ErrCodeMeterModFailedBadRate

	// ErrCodeMeterModFailedBadBurst is returned when specified burst
	// size is unsupported.
	ErrCodeMeterModFailedBadBurst

	// ErrCodeMeterModFailedBadBand is returned when specified band
	// is unsupported.
	ErrCodeMeterModFailedBadBand

	// ErrCodeMeterModFailedBadBandValue is returned when specified
	// band value is unsupported.
	ErrCodeMeterModFailedBadBandValue

	// ErrCodeMeterModFailedOutOfMeters is returned when no more meters
	// available.
	ErrCodeMeterModFailedOutOfMeters

	// ErrCodeMeterModFailedOutOfBands is returned when the maximum
	// number of properties for a meter has been exceeded.
	ErrCodeMeterModFailedOutOfBands
)

const (
	// ErrCodeTableFeaturesFailedBadTable is returned when specified table
	// does not exist.
	ErrCodeTableFeaturesFailedBadTable ErrCode = iota

	// ErrCodeTableFeaturesFailedBadMetadata is returned when specified
	// metadata mask is invalid.
	ErrCodeTableFeaturesFailedBadMetadata

	// ErrCodeTableFeaturesFailedBadType is returned when specified
	// property type is unknown.
	ErrCodeTableFeaturesFailedBadType

	// ErrCodeTableFeaturesFailedBadLen is returned when invalid length
	// was specified in properties.
	ErrCodeTableFeaturesFailedBadLen

	// ErrCodeTableFeaturesFailedBadArgument is returned when unsupported
	// property value was specified.
	ErrCodeTableFeaturesFailedBadArgument

	// ErrCodeTableFeaturesFailedPerm is returned when permission denied.
	ErrCodeTableFeaturesFailedPerm
)

// Error is a message used by the switch to notify the controller of a
// problem.
//
// For example, to create a request to inform the controller about the
// unknown error in the flow modification message:
//
//	req := of.NewRequest(of.TypeError, &Error{
//		Type: ErrTypeFlowModFailed,
//		Code: ErrCodeFlowModFailedUnknown,
//	})
type Error struct {
	// Type value indicates the high-level type of error.
	Type ErrType

	// Code value is interpreted based on the type.
	Code ErrCode

	// Data is variable length and interpreted based on the type and code.
	// Unless specified otherwise, the data field contains at least 64
	// bytes of the failed request that caused the error message to be
	// generated, if the failed request is shorter than 64 bytes it should
	// be the full request without any padding.
	Data []byte
}

func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	errCodeText, ok := errTypeCodeText[e.Type]
	if !ok {
		return fmt.Sprintf("ErrType(%d)Code(%d)", e.Type, e.Code)
	}
	text, ok := errCodeText[e.Code]
	if !ok {
		return fmt.Sprintf("%sCode(%d)", e.Type, e.Code)
	}
	return text
}

var errTypeCodeText = map[ErrType]map[ErrCode]string{
	ErrTypeHelloFailed: {
		ErrCodeHelloFailedIncompatible: "ErrCodeHelloFailedIncompatible",
		ErrCodeHelloFailedPerm:         "ErrCodeHelloFailedPerm",
	},
	ErrTypeBadRequest: {
		ErrCodeBadRequestBadVersion:              "ErrCodeBadRequestBadVersion",
		ErrCodeBadRequestBadType:                 "ErrCodeBadRequestBadType",
		ErrCodeBadRequestBadMultipart:            "ErrCodeBadRequestBadMultipart",
		ErrCodeBadRequestBadExperimenter:         "ErrCodeBadRequestBadExperimenter",
		ErrCodeBadRequestBadExpType:              "ErrCodeBadRequestBadExpType",
		ErrCodeBadRequestPerm:                    "ErrCodeBadRequestPerm",
		ErrCodeBadRequestLen:                     "ErrCodeBadRequestLen",
		ErrCodeBadRequestBufferEmpty:             "ErrCodeBadRequestBufferEmpty",
		ErrCodeBadRequestBufferUnknown:           "ErrCodeBadRequestBufferUnknown",
		ErrCodeBadRequestBadTableID:              "ErrCodeBadRequestBadTableID",
		ErrCodeBadRequestIsSlave:                 "ErrCodeBadRequestIsSlave",
		ErrCodeBadRequestBadPort:                 "ErrCodeBadRequestBadPort",
		ErrCodeBadRequestBadPacket:               "ErrCodeBadRequestBadPacket",
		ErrCodeBadRequestMultipartBufferOverflow: "ErrCodeBadRequestMultipartBufferOverflow",
	},
	ErrTypeBadAction: {
		ErrCodeBadActionType:              "ErrCodeBadActionType",
		ErrCodeBadActionLen:               "ErrCodeBadActionLen",
		ErrCodeBadActionExperimenter:      "ErrCodeBadActionExperimenter",
		ErrCodeBadActionExpType:           "ErrCodeBadActionExpType",
		ErrCodeBadActionOutPort:           "ErrCodeBadActionOutPort",
		ErrCodeBadActionArgument:          "ErrCodeBadActionArgument",
		ErrCodeBadActionPerm:              "ErrCodeBadActionPerm",
		ErrCodeBadActionTooMany:           "ErrCodeBadActionTooMany",
		ErrCodeBadActionQueue:             "ErrCodeBadActionQueue",
		ErrCodeBadActionOutGroup:          "ErrCodeBadActionOutGroup",
		ErrCodeBadActionMatchInconsistent: "ErrCodeBadActionMatchInconsistent",
		ErrCodeBadActionUnsupportedOrder:  "ErrCodeBadActionUnsupportedOrder",
		ErrCodeBadActionTag:               "ErrCodeBadActionTag",
		ErrCodeBadActionSetType:           "ErrCodeBadActionSetType",
		ErrCodeBadActionSetLen:            "ErrCodeBadActionSetLen",
		ErrCodeBadActionSetArgument:       "ErrCodeBadActionSetArgument",
	},
	ErrTypeBadInstruction: {
		ErrCodeBadInstructionUnknown:                 "ErrCodeBadInstructionUnknown",
		ErrCodeBadInstructionUnsupported:             "ErrCodeBadInstructionUnsupported",
		ErrCodeBadInstructionTableID:                 "ErrCodeBadInstructionTableID",
		ErrCodeBadInstructionUnsupportedMetadata:     "ErrCodeBadInstructionUnsupportedMetadata",
		ErrCodeBadInstructionUnsupportedMetadataMask: "ErrCodeBadInstructionUnsupportedMetadataMask",
		ErrCodeBadInstructionExperimenter:            "ErrCodeBadInstructionExperimenter",
		ErrCodeBadInstructionExpType:                 "ErrCodeBadInstructionExpType",
		ErrCodeBadInstructionLen:                     "ErrCodeBadInstructionLen",
		ErrCodeBadInstructionPerm:                    "ErrCodeBadInstructionPerm",
	},
	ErrTypeBadMatch: {
		ErrCodeBadMatchBadType:      "ErrCodeBadMatchBadType",
		ErrCodeBadMatchBadLen:       "ErrCodeBadMatchBadLen",
		ErrCodeBadMatchBadTag:       "ErrCodeBadMatchBadTag",
		ErrCodeBadMatchBadLinkMask:  "ErrCodeBadMatchBadLinkMask",
		ErrCodeBadMatchBadNetMask:   "ErrCodeBadMatchBadNetMask",
		ErrCodeBadMatchBadWildcards: "ErrCodeBadMatchBadWildcards",
		ErrCodeBadMatchBadField:     "ErrCodeBadMatchBadField",
		ErrCodeBadMatchBadValue:     "ErrCodeBadMatchBadValue",
		ErrCodeBadMatchBadMask:      "ErrCodeBadMatchBadMask",
		ErrCodeBadMatchBadPrereq:    "ErrCodeBadMatchBadPrereq",
		ErrCodeBadMatchDupField:     "ErrCodeBadMatchDupField",
		ErrCodeBadMatchPerm:         "ErrCodeBadMatchPerm",
	},
	ErrTypeFlowModFailed: {
		ErrCodeFlowModFailedUnknown:    "ErrCodeFlowModFailedUnknown",
		ErrCodeFlowModFailedTableFull:  "ErrCodeFlowModFailedTableFull",
		ErrCodeFlowModFailedBadTableID: "ErrCodeFlowModFailedBadTableID",
		ErrCodeFlowModFailedOverlap:    "ErrCodeFlowModFailedOverlap",
		ErrCodeFlowModFailedPerm:       "ErrCodeFlowModFailedPerm",
		ErrCodeFlowModFailedBadTimeout: "ErrCodeFlowModFailedBadTimeout",
		ErrCodeFlowModFailedBadCommand: "ErrCodeFlowModFailedBadCommand",
		ErrCodeFlowModFailedBadFlags:   "ErrCodeFlowModFailedBadFlags",
	},
	ErrTypeGroupModFailed: {
		ErrCodeGroupModFailedGroupExists:         "ErrCodeGroupModFailedGroupExists",
		ErrCodeGroupModFailedInvalidGroup:        "ErrCodeGroupModFailedInvalidGroup",
		ErrCodeGroupModFailedWeightUnsupported:   "ErrCodeGroupModFailedWeightUnsupported",
		ErrCodeGroupModFailedOutOfGroups:         "ErrCodeGroupModFailedOutOfGroups",
		ErrCodeGroupModFailedOutOfBuckets:        "ErrCodeGroupModFailedOutOfBuckets",
		ErrCodeGroupModFailedChainingUnsupported: "ErrCodeGroupModFailedChainingUnsupported",
		ErrCodeGroupModFailedWatchUnsupported:    "ErrCodeGroupModFailedWatchUnsupported",
		ErrCodeGroupModFailedLoop:                "ErrCodeGroupModFailedLoop",
		ErrCodeGroupModFailedUnknownGroup:        "ErrCodeGroupModFailedUnknownGroup",
		ErrCodeGroupModFailedChainedGroup:        "ErrCodeGroupModFailedChainedGroup",
		ErrCodeGroupModBadType:                   "ErrCodeGroupModBadType",
		ErrCodeGroupModBadCommand:                "ErrCodeGroupModBadCommand",
		ErrCodeGroupModBadBucket:                 "ErrCodeGroupModBadBucket",
		ErrCodeGroupModBadWatch:                  "ErrCodeGroupModBadWatch",
		ErrCodeGroupModPerm:                      "ErrCodeGroupModPerm",
	},
	ErrTypePortModFailed: {
		ErrCodePortModFailedBadPort:      "ErrCodePortModFailedBadPort",
		ErrCodePortModFailedBadHwAddr:    "ErrCodePortModFailedBadHwAddr",
		ErrCodePortModFailedBadConfig:    "ErrCodePortModFailedBadConfig",
		ErrCodePortModFailedBadAdvertise: "ErrCodePortModFailedBadAdvertise",
		ErrCodePortModFailedPerm:         "ErrCodePortModFailedPerm",
	},
	ErrTypeTableModFailed: {
		ErrCodeTableModFailedBadTable:  "ErrCodeTableModFailedBadTable",
		ErrCodeTableModFailedBadConfig: "ErrCodeTableModFailedBadConfig",
		ErrCodeTableModFailedPerm:      "ErrCodeTableModFailedPerm",
	},
	ErrTypeQueueOpFailed: {
		ErrCodeQueueOpFailedBadPort:  "ErrCodeQueueOpFailedBadPort",
		ErrCodeQueueOpFailedBadQueue: "ErrCodeQueueOpFailedBadQueue",
		ErrCodeQueueOpFailedPerm:     "ErrCodeQueueOpFailedPerm",
	},
	ErrTypeSwitchConfigFailed: {
		ErrCodeSwitchConfigFailedBadFlags: "ErrCodeSwitchConfigFailedBadFlags",
		ErrCodeSwitchConfigFailedBadLen:   "ErrCodeSwitchConfigFailedBadLen",
		ErrCodeSwitchConfigFailedPerm:     "ErrCodeSwitchConfigFailedPerm",
	},
	ErrTypeRoleRequestFailed: {
		ErrCodeRoleRequestFailedStale:   "ErrCodeRoleRequestFailedStale",
		ErrCodeRoleRequestFailedUnsup:   "ErrCodeRoleRequestFailedUnsup",
		ErrCodeRoleRequestFailedBadRole: "ErrCodeRoleRequestFailedBadRole",
	},
	ErrTypeMeterModFailed: {
		ErrCodeMeterModFailedUnknown:      "ErrCodeMeterModFailedUnknown",
		ErrCodeMeterModFailedMeterExists:  "ErrCodeMeterModFailedMeterExists",
		ErrCodeMeterModFailedInvalidMeter: "ErrCodeMeterModFailedInvalidMeter",
		ErrCodeMeterModFailedUnknownMeter: "ErrCodeMeterModFailedUnknownMeter",
		ErrCodeMeterModFailedBadCommand:   "ErrCodeMeterModFailedBadCommand",
		ErrCodeMeterModFailedBadFlags:     "ErrCodeMeterModFailedBadFlags",
		ErrCodeMeterModFailedBadRate:      "ErrCodeMeterModFailedBadRate",
		ErrCodeMeterModFailedBadBurst:     "ErrCodeMeterModFailedBadBurst",
		ErrCodeMeterModFailedBadBand:      "ErrCodeMeterModFailedBadBand",
		ErrCodeMeterModFailedBadBandValue: "ErrCodeMeterModFailedBadBandValue",
		ErrCodeMeterModFailedOutOfMeters:  "ErrCodeMeterModFailedOutOfMeters",
		ErrCodeMeterModFailedOutOfBands:   "ErrCodeMeterModFailedOutOfBands",
	},
	ErrTypeTableFeaturesFailed: {
		ErrCodeTableFeaturesFailedBadTable:    "ErrCodeTableFeaturesFailedBadTable",
		ErrCodeTableFeaturesFailedBadMetadata: "ErrCodeTableFeaturesFailedBadMetadata",
		ErrCodeTableFeaturesFailedBadType:     "ErrCodeTableFeaturesFailedBadType",
		ErrCodeTableFeaturesFailedBadLen:      "ErrCodeTableFeaturesFailedBadLen",
		ErrCodeTableFeaturesFailedBadArgument: "ErrCodeTableFeaturesFailedBadArgument",
		ErrCodeTableFeaturesFailedPerm:        "ErrCodeTableFeaturesFailedPerm",
	},
	ErrTypeExperimenter: {},
}

// WriteTo implements io.WriterTo interface. It serializes the error
// message into the wire format.
func (e *Error) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, e.Type, e.Code, e.Data)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// error message from the wire format.
func (e *Error) ReadFrom(r io.Reader) (n int64, err error) {
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

// ErrorExperimenter defines an experimental error message.
type ErrorExperimenter struct {
	// ExpType is experimenter type defined kind of error.
	ExpType uint16

	// Experimenter identifier.
	Experimenter uint32

	// Data is variable-length error data.
	Data []byte
}

// WriteTo implements io.WriterTo interface. It serializes experimenter
// error message into the wire format.
func (e *ErrorExperimenter) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteTo(w, ErrTypeExperimenter, e.ExpType,
		e.Experimenter, e.Data)
}

// ReadFrom implements io.ReadFrom interface. It deserializes the
// experimenter message from the wire format.
func (e *ErrorExperimenter) ReadFrom(r io.Reader) (n int64, err error) {
	var etype ErrType
	n, err = encoding.ReadFrom(r, &etype, &e.ExpType, &e.Experimenter)
	if err != nil {
		return n, err
	}

	e.Data, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	return n + int64(len(e.Data)), nil
}
