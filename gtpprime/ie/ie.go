// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package ie provides encoding/decoding feature of gtpprime Information Elements.
*/
package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// IE definitions.
const (
	Cause                             uint8 = 2
	PacketTransferCommand             uint8 = 126
	SequenceNumbersOfReleasedPackets  uint8 = 249
	SequenceNumbersOfCancelledPackets uint8 = 250
	DataRecordPacket                  uint8 = 252
	RequestsResponded                 uint8 = 253
	PrivateExtension                  uint8 = 255
)

// IE is a gtpprime Information Element.
type IE struct {
	Type     uint8
	Length   uint16
	Payload  []byte
	ChildIEs []*IE
}

// New creates new IE.
func New(itype, ins uint8, data []byte) *IE {
	ie := &IE{
		Type:    itype,
		Payload: data,
	}
	ie.SetLength()

	return ie
}

// Marshal returns the byte sequence generated from an IE instance.
func (i *IE) Marshal() ([]byte, error) {
	b := make([]byte, i.MarshalLen())
	if err := i.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (i *IE) MarshalTo(b []byte) error {
	l := len(b)
	requiredLength := 1
	if i.HasLength() {
		requiredLength = 3
	}
	if l < requiredLength {
		return io.ErrUnexpectedEOF
	}

	b[0] = i.Type
	offset := 1
	if i.HasLength() {
		binary.BigEndian.PutUint16(b[offset:offset+2], i.Length)
		offset += 2
	}

	//if i.IsGrouped() {
	//	offset += 2
	//	for _, ie := range i.ChildIEs {
	//		if l < offset+ie.MarshalLen() {
	//			break
	//		}
	//
	//		if err := ie.MarshalTo(b[offset : offset+ie.MarshalLen()]); err != nil {
	//			return err
	//		}
	//		offset += ie.MarshalLen()
	//	}
	//	return nil
	//}

	if l < i.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	copy(b[offset:i.MarshalLen()], i.Payload)
	return nil
}

// Parse decodes given byte sequence as a gtpprime Information Element.
func Parse(b []byte) (*IE, error) {
	ie := &IE{}
	if err := ie.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return ie, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in gtpprime IE.
func (i *IE) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 4 {
		return io.ErrUnexpectedEOF
	}

	i.Type = b[0]
	i.Length = binary.BigEndian.Uint16(b[1:3])
	if int(i.Length) > l-4 {
		return ErrInvalidLength
	}

	//	i.instance = b[3]
	i.Payload = b[4 : 4+int(i.Length)]

	if i.IsGrouped() {
		var err error
		i.ChildIEs, err = ParseMultiIEs(i.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalLen returns field length in integer.
func (i *IE) MarshalLen() int {
	l := 1
	if i.HasLength() {
		l += 2
	}
	//if i.IsGrouped() {
	//	l := 4
	//	for _, ie := range i.ChildIEs {
	//		l += ie.MarshalLen()
	//	}
	//	return l
	//}
	return l + len(i.Payload)
}

// SetLength sets the length in Length field.
func (i *IE) SetLength() {
	if i.IsGrouped() {
		l := 0
		for _, ie := range i.ChildIEs {
			l += ie.MarshalLen()
		}
		i.Length = uint16(l)
	}
	i.Length = uint16(len(i.Payload))
}

// Name returns the name of IE in string.
func (i *IE) Name() string {
	if n, ok := ieTypeNameMap[i.Type]; ok {
		return n
	}
	return "Undefined"
}

// String returns the gtpprime IE values in human readable format.
func (i *IE) String() string {
	if i == nil {
		return "nil"
	}
	return fmt.Sprintf("{%s: {Type: %d, Length: %d, Payload: %#v}}",
		i.Name(),
		i.Type,
		i.Length,
		i.Payload,
	)
}

var grouped = []uint8{}

// IsGrouped reports whether an IE is grouped type or not.
func (i *IE) IsGrouped() bool {
	for _, itype := range grouped {
		if i.Type == itype {
			return true
		}
	}
	return false
}

// HasLength reports whether an IE has length or not.
func (i *IE) HasLength() bool {
	if (i.Type & 128) < 128 {
		return false
	}
	return true
}

// Add adds variable number of IEs to a IE if the IE is grouped type and update length.
// Otherwise, this does nothing(no errors).
func (i *IE) Add(ies ...*IE) {
	if !i.IsGrouped() {
		return
	}

	i.Payload = nil
	i.ChildIEs = append(i.ChildIEs, ies...)
	for _, ie := range i.ChildIEs {
		serialized, err := ie.Marshal()
		if err != nil {
			continue
		}
		i.Payload = append(i.Payload, serialized...)
	}
	i.SetLength()
}

// Remove removes an IE looked up by type and instance.
func (i *IE) Remove(typ, instance uint8) {
	if !i.IsGrouped() {
		return
	}

	i.Payload = nil
	newChildren := make([]*IE, len(i.ChildIEs))
	idx := 0
	for _, ie := range i.ChildIEs {
		if ie.Type == typ {
			newChildren = newChildren[:len(newChildren)-1]
			continue
		}
		newChildren[idx] = ie
		idx++

		serialized, err := ie.Marshal()
		if err != nil {
			continue
		}
		i.Payload = append(i.Payload, serialized...)
	}
	i.ChildIEs = newChildren
	i.SetLength()
}

// FindByType returns IE looked up by type and instance.
//
// The program may be slower when calling this method multiple times
// because this ranges over a ChildIEs each time it is called.
func (i *IE) FindByType(typ, instance uint8) (*IE, error) {
	if !i.IsGrouped() {
		return nil, ErrInvalidType
	}

	for _, ie := range i.ChildIEs {
		if ie.Type == typ {
			return ie, nil
		}
	}
	return nil, ErrIENotFound
}

// ParseMultiIEs decodes multiple IEs at a time.
// This is easy and useful but slower than decoding one by one.
// When you don't know the number of IEs, this is the only way to decode them.
// See benchmarks in diameter_test.go for the detail.
func ParseMultiIEs(b []byte) ([]*IE, error) {
	var ies []*IE
	for {
		if len(b) == 0 {
			break
		}

		i, err := Parse(b)
		if err != nil {
			return nil, err
		}
		ies = append(ies, i)
		b = b[i.MarshalLen():]
	}
	return ies, nil
}

func newUint8ValIE(t, v uint8) *IE {
	return New(t, 0x00, []byte{v})
}

func newUint16ValIE(t uint8, v uint16) *IE {
	i := New(t, 0x00, make([]byte, 2))
	binary.BigEndian.PutUint16(i.Payload, v)
	return i
}

func newUint32ValIE(t uint8, v uint32) *IE {
	i := New(t, 0x00, make([]byte, 4))
	binary.BigEndian.PutUint32(i.Payload, v)
	return i
}

// unused for now.
// func newUint64ValIE(t uint8, v uint64) *IE {
// 	i := New(t, 0x00, make([]byte, 8))
// 	binary.BigEndian.PutUint64(i.Payload, v)
// 	return i
// }

func newStringIE(t uint8, v string) *IE {
	return New(t, 0x00, []byte(v))
}

func newGroupedIE(itype uint8, ies ...*IE) *IE {
	i := New(itype, 0x00, make([]byte, 0))
	i.ChildIEs = ies
	for _, ie := range i.ChildIEs {
		serialized, err := ie.Marshal()
		if err != nil {
			return nil
		}
		i.Payload = append(i.Payload, serialized...)
	}
	i.SetLength()

	return i
}

var ieTypeNameMap = map[uint8]string{
	126: "PacketTransferCommand",
}
