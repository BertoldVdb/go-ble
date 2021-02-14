package bleatt

import (
	"errors"
)

var (
	ErrorProtocolViolation = errors.New("Protocal violation detected")
)

type ATTCommand uint8

const (
	ATTErrorRsp               ATTCommand = 0x1
	ATTExchangeMTUReq         ATTCommand = 0x2
	ATTExchangeMTURsp         ATTCommand = 0x3
	ATTFindInformationReq     ATTCommand = 0x4
	ATTFindInformationRsp     ATTCommand = 0x5
	ATTFindByTypeValueReq     ATTCommand = 0x6
	ATTFindByTypeValueRsp     ATTCommand = 0x7
	ATTReadByTypeReq          ATTCommand = 0x8
	ATTReadByTypeRsp          ATTCommand = 0x9
	ATTReadReq                ATTCommand = 0xA
	ATTReadRsp                ATTCommand = 0xB
	ATTReadBlobReq            ATTCommand = 0xC
	ATTReadBlobRsp            ATTCommand = 0xD
	ATTReadMultipleReq        ATTCommand = 0xE
	ATTReadMultipleRsp        ATTCommand = 0xF
	ATTReadByGroupTypeReq     ATTCommand = 0x10
	ATTReadByGroupTypeRsp     ATTCommand = 0x11
	ATTWriteReq               ATTCommand = 0x12
	ATTWriteRsp               ATTCommand = 0x13
	ATTWriteCMD               ATTCommand = 0x52
	ATTPrepareWriteReq        ATTCommand = 0x16
	ATTPrepareWriteRsp        ATTCommand = 0x17
	ATTExecuteWriteReq        ATTCommand = 0x18
	ATTExecuteWriteRsp        ATTCommand = 0x19
	ATTReadMultipleValueReq   ATTCommand = 0x20
	ATTReadMultipleValueRsp   ATTCommand = 0x21
	ATTMultipleHandleValueNTF ATTCommand = 0x23
	ATTHandleValueNTF         ATTCommand = 0x1B
	ATTHandleValueIND         ATTCommand = 0x1D
	ATTHandleValueCNF         ATTCommand = 0x1E
	ATTSignedWriteCMD         ATTCommand = 0xD2
)

type ATTError uint8

const (
	ATTErrorNone                       ATTError = 0
	ATTErrorInvalidHandle              ATTError = 0x01
	ATTErrorReadNotPermitted           ATTError = 0x02
	ATTErrorWriteNotPermitted          ATTError = 0x03
	ATTErrorInvalidPDU                 ATTError = 0x04
	ATTErrorInsufficientAuthentication ATTError = 0x05
	ATTErrorRequestNotSupported        ATTError = 0x06
	ATTErrorInvalidOffset              ATTError = 0x07
	ATTErrorInsufficientAuthorization  ATTError = 0x08
	ATTErrorPrepareQueueFull           ATTError = 0x09
	ATTErrorAttributeNotFound          ATTError = 0x0A
	ATTErrorAttributeNotLong           ATTError = 0x0B
	ATTErrorSize                       ATTError = 0x0C
	ATTErrorLength                     ATTError = 0x0D
	ATTErrorUnlikelyError              ATTError = 0x0E
	ATTErrorInsufficientEncryption     ATTError = 0x0F
	ATTErrorUnsupportedGroupType       ATTError = 0x10
	ATTErrorInsufficientResources      ATTError = 0x11
	ATTErrorDatabaseOutOfSync          ATTError = 0x12
	ATTErrorValueNotAllowed            ATTError = 0x13
)
