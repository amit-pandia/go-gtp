// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"github.com/wmnsk/go-gtp/gtpv2/testutils"
)

func TestReleaseAccessBearersResponse(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal/CauseOnly",
			Structured: message.NewReleaseAccessBearersResponse(
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
			),
			Serialized: []byte{
				// Header
				0x48, 0xab, 0x00, 0x0e, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00, 0x01, 0x00,
				0x02, 0x00, 0x02, 0x00, 0x10, 0x00,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseReleaseAccessBearersResponse(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
