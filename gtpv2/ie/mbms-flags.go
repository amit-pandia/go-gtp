// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewMBMSFlags creates a new MBMSFlags IE.
func NewMBMSFlags(lmri, msri uint8) *IE {
	i := New(MBMSFlags, 0x00, make([]byte, 1))
	i.Payload[0] |= (lmri << 1 & 0x02) | (msri & 0x01)
	return i
}

// MBMSFlags returns MBMSFlags in uint8 if the type of IE matches.
func (i *IE) MBMSFlags() (uint8, error) {
	if i.Type != MBMSFlags {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[0], nil
}

// MustMBMSFlags returns MBMSFlags in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustMBMSFlags() uint8 {
	v, _ := i.MBMSFlags()
	return v
}

// HasMSRI reports whether an IE has MSRI bit.
func (i *IE) HasMSRI() bool {
	v, err := i.MBMSFlags()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasLMRI reports whether an IE has LMRI bit.
func (i *IE) HasLMRI() bool {
	v, err := i.MBMSFlags()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// LocalMBMSBearerContextRelease reports whether the MBMS Session Stop Request
// message is used to release the MBMS Bearer Context locally in the MME/SGSN.
func (i *IE) LocalMBMSBearerContextRelease() bool {
	if len(i.Payload) < 1 {
		return false
	}
	switch i.Type {
	case MBMSFlags:
		return i.Payload[0]&0x02 == 1
	default:
		return false
	}
}

// MBMSSessionReEstablishment reports whether the MBMS Session Start Request
// message is used to re-establish an MBMS session.
func (i *IE) MBMSSessionReEstablishment() bool {
	if len(i.Payload) < 1 {
		return false
	}
	switch i.Type {
	case MBMSFlags:
		return i.Payload[0]&0x01 == 1
	default:
		return false
	}
}
