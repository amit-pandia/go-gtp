// Copyright 2019 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ies

import "github.com/pkg/errors"

// Error definitions.
var (
	ErrInvalidLength       = errors.New("got invalid length ")
	ErrTooShortToMarshal = errors.New("too short to serialize")
	ErrTooShortToParse    = errors.New("too short to decode as GTPv1 IE")
)
