// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"github.com/wmnsk/go-gtp/gtpv2/testutils"
)

func TestContextRequest(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewContextRequest(
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewIMSI("123451234567890"),
				ie.NewPacketTMSI(0xdeadbeef),
				ie.NewPTMSISignature(0xbeebee),
			),
			Serialized: []byte{
				// Header
				0x48, 0x82, 0x00, 0x23, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00, 0x01, 0x00,
				// IMSI
				0x01, 0x00, 0x08, 0x00, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				// P-TMSI
				0x6f, 0x00, 0x04, 0x00, 0xde, 0xad, 0xbe, 0xef,
				// P-TMSI Signature
				0x70, 0x00, 0x03, 0x00, 0xbe, 0xeb, 0xee,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseContextRequest(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
