// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smos

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSmoS_Balance .
func TestSmoS_Balance(t *testing.T) {
	smos := NewSmoS()
	smos.Add("192.168.1.1", 0)
	smos.Add("192.168.1.2", 80)
	smos.Remove("192.168.1.2")
	smos.Inc("192.168.1.1")
	smos.Inc("192.168.1.1")
	smos.Done("192.168.1.1")
	host, err := smos.Balance(&Function{
		FuncUUID:     "xj92-3242-csxa-JKjx",
		FuncName:     "hello world",
		LibSortedSet: []string{"testify", "quic"},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, host, "192.168.1.1")
	smos.Remove("192.168.1.1")
	host, err = smos.Balance(&Function{
		FuncUUID:     "xj92-3242-csxa-JKjx",
		FuncName:     "hello world",
		LibSortedSet: []string{"testify", "quic"},
	})
	assert.Equal(t, err, errors.New("no hosts added"))
	assert.Equal(t, host, "")
}
