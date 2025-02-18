// Copyright 2019-2022 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"errors"
	"fmt"
)

// Error definitions.
var (
	ErrInvalidLength     = errors.New("got invalid length ")
	ErrTooShortToMarshal = errors.New("too short to serialize")
	ErrTooShortToParse   = errors.New("too short to decode as GTPv1 IE")

	ErrMalformed = errors.New("malformed IE")
)

// InvalidTypeError indicates the type of IE is invalid.
type InvalidTypeError struct {
	Type uint8
}

// Error returns message with the invalid type given.
func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("got invalid type: %v", e.Type)
}
