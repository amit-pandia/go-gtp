// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/amit-pandia/go-gtp/gtpprime/ie"
)

// DataRecordTransferRequest is a DataRecordTransferRequest Header and its IEs above.
type DataRecordTransferRequest struct {
	*Header
	PacketTransferCommand *ie.IE
}

// NewDataRecordTransferRequest creates a new DataRecordTransferRequest.
func NewDataRecordTransferRequest(seq uint16, ies ...*ie.IE) *DataRecordTransferRequest {
	e := &DataRecordTransferRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 0),
			MsgTypeDataRecordTransferRequest, seq, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PacketTransferCommand:
			e.PacketTransferCommand = i
		}
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a DataRecordTransferRequest.
func (e *DataRecordTransferRequest) Marshal() ([]byte, error) {
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *DataRecordTransferRequest) MarshalTo(b []byte) error {
	if e.Header.Payload != nil {
		e.Header.Payload = nil
	}
	e.Header.Payload = make([]byte, e.MarshalLen()-e.Header.MarshalLen())

	offset := 0
	if ie := e.PacketTransferCommand; ie != nil {
		if err := ie.MarshalTo(e.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += ie.MarshalLen()
	}
	e.Header.SetLength()
	return e.Header.MarshalTo(b)
}

// ParseDataRecordTransferRequest decodes a given byte sequence as a DataRecordTransferRequest.
func ParseDataRecordTransferRequest(b []byte) (*DataRecordTransferRequest, error) {
	e := &DataRecordTransferRequest{}
	if err := e.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a DataRecordTransferRequest.
func (e *DataRecordTransferRequest) UnmarshalBinary(b []byte) error {
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if len(e.Header.Payload) < 2 {
		return nil
	}

	decodedIEs, err := ie.ParseMultiIEs(e.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range decodedIEs {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PacketTransferCommand:
			e.PacketTransferCommand = i
		}
	}

	return nil
}

// MarshalLen returns the serial length of Data.
func (e *DataRecordTransferRequest) MarshalLen() int {
	l := e.Header.MarshalLen() - len(e.Header.Payload)

	if ie := e.PacketTransferCommand; ie != nil {
		l += ie.MarshalLen()
	}

	return l
}

// SetLength sets the length in Length field.
func (e *DataRecordTransferRequest) SetLength() {
	e.Header.Length = uint16(e.MarshalLen() - 4)
}

// MessageTypeName returns the name of protocol.
func (e *DataRecordTransferRequest) MessageTypeName() string {
	return "Echo Request"
}
