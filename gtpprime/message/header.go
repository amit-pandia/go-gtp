// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"fmt"
)

const (
	fixedHeaderSize  = 4
	seqSpareSize     = 4
	teidSize         = 4
	noTEIDHeaderSize = fixedHeaderSize + seqSpareSize
	teidHeaderSize   = noTEIDHeaderSize + teidSize
)

// Header is a gtpprime common header
type Header struct {
	Flags          uint8
	Type           uint8
	Length         uint16
	SequenceNumber uint16
	Payload        []byte
}

// NewHeader creates a new Header
func NewHeader(flags, mtype uint8, seqnum uint16, data []byte) *Header {
	h := &Header{
		Flags:          flags,
		Type:           mtype,
		SequenceNumber: seqnum,
		Payload:        data,
	}
	h.SetLength()
	fmt.Println("gtpp hdr len =", h.Length)

	return h
}

// NewHeaderFlags returns a Header Flag built by its components given as arguments.
func NewHeaderFlags(v, p, s int) uint8 {
	return uint8(
		((v & 0x7) << 5) | ((p & 0x1) << 4) | ((s & 0x7) << 1),
	)
}

// Marshal returns the byte sequence generated from a Header instance.
func (h *Header) Marshal() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Header) MarshalTo(b []byte) error {
	b[0] = h.Flags
	b[1] = h.Type
	binary.BigEndian.PutUint16(b[2:4], h.Length)
	offset := 4
	binary.BigEndian.PutUint16(b[offset:offset+2], h.SequenceNumber)
	copy(b[offset+2:h.MarshalLen()], h.Payload)
	fmt.Println("b=", b)

	return nil
}

// ParseHeader decodes given byte sequence as a gtpprime header.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in gtpprime header.
func (h *Header) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 6 {
		return ErrTooShortToParse
	}
	h.Flags = b[0]
	h.Type = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:4])
	if h.Length < seqSpareSize {
		return ErrTooShortToParse
	}
	h.Length = binary.BigEndian.Uint16(b[4:6])

	if int(h.Length)+fixedHeaderSize > l {
		h.Payload = b[noTEIDHeaderSize:]
		return nil
	}
	if fixedHeaderSize+h.Length >= noTEIDHeaderSize {
		h.Payload = b[noTEIDHeaderSize : fixedHeaderSize+h.Length]
	} else {
		return ErrInvalidLength
	}

	return nil
}

// MarshalLen returns field length in integer.
func (h *Header) MarshalLen() int {
	l := 6 + len(h.Payload)
	return l
}

// SetLength sets the length in Length field.
func (h *Header) SetLength() {
	h.Length = uint16(4 + len(h.Payload))
	if h.HasTEID() {
		h.Length += 4
	}
}

// String returns the gtpprime header values in human readable format.
func (h *Header) String() string {
	return fmt.Sprintf("{Flags: %#x, Type: %d, Length: %d, SequenceNumber: %#x, Payload: %#v}",
		h.Flags,
		h.Type,
		h.Length,
		h.SequenceNumber,
		h.Payload,
	)
}

// IsPiggybacking reports whether the message has the trailing(piggybacked) message.
func (h *Header) IsPiggybacking() bool {
	return (int(h.Flags)>>4)&0x01 == 1
}

// SetPiggybacking sets the Piggybacking flag.
//
// The given value should only be 0 or 1. Otherwise it may cause the unexpected result.
func (h *Header) SetPiggybacking(val uint8) {
	h.Flags = (h.Flags & 0xef) | (val & 0x01 << 4)
}

// HasTEID determines whether a gtpprime has TEID inside by checking the flag.
func (h *Header) HasTEID() bool {
	return (int(h.Flags)>>3)&0x01 == 1
}

// Sequence returns SequenceNumber in uint16.
func (h *Header) Sequence() uint16 {
	return h.SequenceNumber
}

// SetSequenceNumber sets the SequenceNumber in Header.
func (h *Header) SetSequenceNumber(seq uint16) {
	h.SequenceNumber = seq
}

// HasMessagePriority reports whether the message has MessagePriority field
func (h *Header) HasMessagePriority() bool {
	return (int(h.Flags)>>2)&0x01 == 1
}

// SetMessagePriority sets the MessagePriorityFlag to 1 and puts the MessagePriority
// given into MessagePriority field.
func (h *Header) SetMessagePriority(mp uint8) {
	h.Flags |= (1 << 2)
}

// Version returns the GTP version.
func (h *Header) Version() int {
	return 2
}

// MessageType returns the type of messagg.
func (h *Header) MessageType() uint8 {
	return h.Type
}
