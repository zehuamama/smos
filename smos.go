// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smos

import "sync"

const salt = "%$#"

// Function represents a function, including its UUID,
// function name, and library file of the function
type Function struct {
	FuncUUID      string
	FuncName      string
	FuncSortedSet []string
}

// SmoS refers to serverless multi-objective scheduling algorithm
type SmoS struct {
	sync.RWMutex
	c          *Consistent
	maxLoadMap map[string]uint64
}

// NewSmoS ...
func NewSmoS() *SmoS {
	return &SmoS{
		c:          NewConsistent(),
		maxLoadMap: make(map[string]uint64),
	}
}

// Add host to the ring
func (s *SmoS) Add(host string, maxLoad uint64) {
	s.Lock()
	defer s.Unlock()
	s.maxLoadMap[host] = maxLoad
	s.c.Add(host)
}

// Remove host from the ring
func (s *SmoS) Remove(host string, maxLoad uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.maxLoadMap, host)
	s.c.Remove(host)
}

// Inc host load increases by 1
func (s *SmoS) Inc(host string) {
	s.Lock()
	defer s.Unlock()
	s.c.Inc(host)
}

// Done host load reduced by 1
func (s *SmoS) Done(host string) {
	s.Lock()
	defer s.Unlock()
	s.c.Done(host)
}

// Balance chooses the appropriate host to call function
func (s *SmoS) Balance(function Function) (string, error) {
	s.RLock()
	defer s.RUnlock()
	assign := ""
	if len(function.FuncSortedSet) != 0 {
		node1, err := s.c.Get(function.FuncSortedSet[0])
		if err != nil {
			return "", err
		}
		node2, err := s.c.Get(function.FuncSortedSet[0] + salt)
		if err != nil {
			return "", err
		}
		if s.c.GetLoad(node1) < s.c.GetLoad(node2) {
			assign = node1
		} else {
			assign = node2
		}
		if s.c.GetLoad(assign) < s.maxLoadMap[assign] {
			return assign, nil
		}
	}
	assign, err := s.c.GetLeast(function.FuncUUID)
	if err != nil {
		return "", err
	}
	return assign, nil
}
