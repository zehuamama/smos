// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smos

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// TestConsistent_Add .
func TestConsistent_Add(t *testing.T) {
	c := NewConsistent()
	c.Add("172.166.1.1")
	c.Add("172.166.1.2")
	c.Add("172.166.1.2")
	expect := NewConsistent()
	expect.Add("172.166.1.2")
	expect.Add("172.166.1.1")
	assert.Equal(t, true, reflect.DeepEqual(c, expect))
}

// TestConsistent_Remove .
func TestConsistent_Remove(t *testing.T) {
	c := NewConsistent()
	c.Add("172.166.1.1")
	c.Add("172.166.1.2")
	c.Remove("172.166.1.2")
	expect := NewConsistent()
	expect.Add("172.166.1.1")
	assert.Equal(t, true, reflect.DeepEqual(c, expect))
}

// TestConsistent_GetLoad .
func TestConsistent_GetLoad(t *testing.T) {
	c := NewConsistent()
	c.Add("172.166.1.1")
	c.Add("172.166.1.2")
	assert.Equal(t, true, reflect.DeepEqual(c.GetLoad("172.166.1.1"), uint64(0)))
}

// TestConsistent_Get .
func TestConsistent_Get(t *testing.T) {
	c := NewConsistent()
	c.Add("172.166.1.1")
	c.Add("172.166.1.2")
	c.Inc("172.166.1.1")
	c.Inc("nil")
	c.Done("172.166.1.1")
	c.Done("nil")
	host, err := c.Get("192.168.0.1")
	assert.Equal(t, true, reflect.DeepEqual(err, nil))
	assert.Equal(t, true, reflect.DeepEqual(host, "172.166.1.1"))
}

// TestConsistent_GetLeast .
func TestConsistent_GetLeast(t *testing.T) {
	c := NewConsistent()
	c.Add("172.166.1.1")
	c.Add("172.166.1.2")
	host, err := c.GetLeast("192.168.0.1")
	assert.Equal(t, true, reflect.DeepEqual(err, nil))
	assert.Equal(t, true, reflect.DeepEqual(host, "172.166.1.1"))
	c.Inc("172.166.1.1")
	c.Inc("172.166.1.1")
	host, err = c.GetLeast("192.168.0.1")
	assert.Equal(t, true, reflect.DeepEqual(err, nil))
	assert.Equal(t, true, reflect.DeepEqual(host, "172.166.1.2"))
}
