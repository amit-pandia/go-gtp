// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package message provides encoding/decoding feature of gtpprime protocol.
*/
package message

import (
	"fmt"
	"github.com/amit-pandia/go-gtp/gtpprime/ie"
	"io"
	"reflect"
	"strings"
)

// Message Type definitions.
const (
	_ uint8 = iota
	MsgTypeEchoRequest
	MsgTypeEchoResponse
	MsgTypeVersionNotSupportedIndication
	_
	_
	_
	_
	_ uint8 = iota + 240 - 9
	MsgTypeDataRecordTransferRequest
	MsgTypeDataRecordTransferResponse
)

// Message is an interface that defines gtpprime message.
type Message interface {
	MarshalTo([]byte) error
	UnmarshalBinary(b []byte) error
	MarshalLen() int
	Version() int
	MessageType() uint8
	MessageTypeName() string
	Sequence() uint16
	SetSequenceNumber(uint16)
}

// Marshal returns the byte sequence generated from a Message instance.
// Better to use MarshalXxx instead if you know the name of message to be serialized.
func Marshal(m Message) ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Parse decodes the given bytes as Message.
func Parse(b []byte) (Message, error) {
	var m Message

	if len(b) < 2 {
		return nil, io.ErrUnexpectedEOF
	}

	switch b[1] {
	case MsgTypeDataRecordTransferRequest:
		m = &DataRecordTransferRequest{}
	default:
		m = &Generic{}
	}

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, fmt.Errorf("failed to decode gtpprime Message: %w", err)
	}
	return m, nil
}

// Prettify returns a Message in prettified representation in string.
//
// Note that this relies much on reflect package, and thus the frequent use of
// this function may have a serious impact on the performance of your software.
func Prettify(m Message) string {
	name := m.MessageTypeName()
	header := strings.TrimSuffix(fmt.Sprint(m), "}")

	v := reflect.Indirect(reflect.ValueOf(m))
	n := v.NumField() - 1
	fields := make([]*field, n)
	for i := 1; i < n+1; i++ { // Skip *Header
		fields[i-1] = &field{name: v.Type().Field(i).Name, maybeIE: v.Field(i).Interface()}
	}

	return fmt.Sprintf("{%s: %s, IEs: [%v]}", name, header, strings.Join(prettifyFields(fields), ", "))
}

type field struct {
	name    string
	maybeIE interface{}
}

func prettifyFields(fields []*field) []string {
	ret := []string{}
	for _, field := range fields {
		if field.maybeIE == nil {
			ret = append(ret, prettifyIE(field.name, nil))
			continue
		}

		// TODO: do this recursively?
		v, ok := field.maybeIE.(*ie.IE)
		if !ok {
			// only for AdditionalIEs field
			if ies, ok := field.maybeIE.([]*ie.IE); ok {
				vals := make([]string, len(ies))
				for i, val := range ies {
					vals[i] = fmt.Sprint(val)
				}
				ret = append(ret, fmt.Sprintf("{%s: [%v]}", field.name, strings.Join(vals, ", ")))
			}
			continue
		}

		ret = append(ret, prettifyIE(field.name, v))
	}

	return ret
}

func prettifyIE(name string, i *ie.IE) string {
	if i == nil {
		return fmt.Sprintf("{%s: %v}", name, i)
	}

	if i.IsGrouped() {
		vals := make([]string, len(i.ChildIEs))
		for i, val := range i.ChildIEs {
			vals[i] = fmt.Sprint(val)
		}
		return fmt.Sprintf("{%s: [%v]}", name, strings.Join(vals, ", "))
	}

	return fmt.Sprintf("{%s: %v}", name, i)
}
