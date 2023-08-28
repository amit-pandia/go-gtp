// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/amit-pandia/go-gtp/gtpprime/ie"
	"github.com/amit-pandia/go-gtp/gtpprime/message"
	"github.com/amit-pandia/go-gtp/gtpprime/testutils"
)

func TestDataRecordTransferRequest(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewDataRecordTransferRequest(
				5,
				ie.NewPacketTransferCommand(1),
			),
			Serialized: []byte{
				0x4e, 0xf0, 0x01, 0xfa, 0x67, 0x2b, 0x7e, 0x01,
				// Packet Transfer Command Type
				0x7e,
				// Packet Transfer Command
				0x01,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseDataRecordTransferRequest(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
