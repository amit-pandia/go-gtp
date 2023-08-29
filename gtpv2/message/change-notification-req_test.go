// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/amit-pandia/go-gtp/gtpv2"
	"github.com/amit-pandia/go-gtp/gtpv2/ie"
	"github.com/amit-pandia/go-gtp/gtpv2/message"
	"github.com/amit-pandia/go-gtp/gtpv2/testutils"
)

func TestChangeNotificationRequest(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewChangeNotificationRequest(
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewIMSI("123451234567890"),
				ie.NewMobileEquipmentIdentity("123450123456789"),
				ie.NewIndication(
					1, 0, 1, 0, 0, 0, 0, 1,
					0, 0, 0, 0, 1, 0, 0, 0,
					0, 0, 0, 1, 0, 1, 0, 1,
					0, 0, 0, 1, 0, 0, 0, 0,
					1, 0, 0, 0, 1, 0, 0, 0,
					1, 0, 0, 0, 0, 0, 0, 1,
					0, 1, 0, 0, 0, 0, 0, 0,
					1, 0, 1, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 1,
				),
				ie.NewRATType(gtpv2.RATTypeEUTRAN),
				ie.NewUserLocationInformationStruct(
					ie.NewCGI("123", "45", 0x1111, 0x2222),
					ie.NewSAI("123", "45", 0x1111, 0x3333),
					ie.NewRAI("123", "45", 0x1111, 0x4444),
					ie.NewTAI("123", "45", 0x5555),
					ie.NewECGI("123", "45", 0x66666666),
					ie.NewLAI("123", "45", 0x1111),
					ie.NewMENBI("123", "45", 0x11111111),
					ie.NewEMENBI("123", "45", 0x22222222),
				),
				ie.NewUserCSGInformation("123", "45", 0x00ffffff, gtpv2.AccessModeHybrid, 0, gtpv2.CMICSG),
				ie.NewIPAddress("1.1.1.1"),
				ie.NewEPSBearerID(0x05),
				ie.NewPrivateExtension(10415, []byte{0xde, 0xad, 0xbe, 0xef}),
			),
			Serialized: []byte{
				// Header
				0x48, 0x26, 0x00, 0x8c, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00, 0x01, 0x00,
				// IMSI
				0x01, 0x00, 0x08, 0x00, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				// MAI
				0x4b, 0x00, 0x08, 0x00, 0x21, 0x43, 0x05, 0x21, 0x43, 0x65, 0x87, 0xf9,
				// IndicationFlags
				0x4d, 0x00, 0x09, 0x00, 0xa1, 0x08, 0x15, 0x10, 0x88, 0x81, 0x40, 0xa0, 0x01,
				// RATType
				0x52, 0x00, 0x01, 0x00, 0x06,
				// ULI
				// --------------------------------------
				0x56, 0x00, 0x33, 0x00,
				// Flags
				0xff,
				// CGI
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x22, 0x22,
				// SAI
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x33, 0x33,
				// RAI
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x44, 0x44,
				// TAI
				0x21, 0xf3, 0x54, 0x55, 0x55,
				// ECGI
				0x21, 0xf3, 0x54, 0x06, 0x66, 0x66, 0x66,
				// RAI
				0x21, 0xf3, 0x54, 0x11, 0x11,
				// Macro eNB ID
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x11,
				// Extended Macro eNB ID
				0x21, 0xf3, 0x54, 0x22, 0x22, 0x22,
				// --------------------------------------
				// UCI
				0x91, 0x00, 0x08, 0x00, 0x21, 0xf3, 0x54, 0x00, 0xff, 0xff, 0xff, 0x41,
				// PGWS5S8IPAddressForControlPlane
				0x4a, 0x00, 0x04, 0x00, 0x01, 0x01, 0x01, 0x01,
				// LinkedEBI
				0x49, 0x00, 0x01, 0x00, 0x05,
				// PrivateExtension
				0xff, 0x00, 0x06, 0x00, 0x28, 0xaf, 0xde, 0xad, 0xbe, 0xef,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseChangeNotificationRequest(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
