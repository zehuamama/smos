// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smos

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

const (
	replicationFactor = 10
	loadFactor        = 1.25
)

var (
	ErrNoHosts = errors.New("no hosts added")
)

// Host from the ring
type Host struct {
	Name string
	Load uint64
}

// Consistent refers to consistent hash with bounded
type Consistent struct {
	hosts     map[uint32]string
	sortedSet []uint32
	loadMap   map[string]*Host
	totalLoad uint64
}

// NewConsistent ...
func NewConsistent() *Consistent {
	return &Consistent{
		hosts:     make(map[uint32]string),
		loadMap:   make(map[string]*Host),
		sortedSet: make([]uint32, 0),
	}
}

// Add host to the ring
func (c *Consistent) Add(host string) {
	if _, ok := c.loadMap[host]; ok {
		return
	}

	c.loadMap[host] = &Host{Name: host, Load: 0}
	for i := 0; i < replicationFactor; i++ {
		hostName := fmt.Sprintf("%s%d", host, i)
		h := hash([]byte(hostName))
		c.hosts[h] = host
		c.sortedSet = append(c.sortedSet, h)

	}
	sort.Slice(c.sortedSet, func(i int, j int) bool {
		return c.sortedSet[i] < c.sortedSet[j]
	})
}

// Remove host from the ring
func (c *Consistent) Remove(host string) bool {
	for i := 0; i < replicationFactor; i++ {
		hostName := fmt.Sprintf("%s%d", host, i)
		h := hash([]byte(hostName))
		delete(c.hosts, h)
		deleteSlice(&c.sortedSet, h)
	}
	delete(c.loadMap, host)
	return true
}

// Get chooses the appropriate host (according to the consistent hash algorithm)
func (c *Consistent) Get(key string) (string, error) {
	if len(c.hosts) == 0 {
		return "", ErrNoHosts
	}

	h := hash([]byte(key))
	idx := searchSlice(c.sortedSet, h)
	return c.hosts[c.sortedSet[idx]], nil
}

// GetLeast chooses the appropriate host (according to the consistent hash with bounded)
// refers to https://research.googleblog.com/2017/04/consistent-hashing-with-bounded-loads.html
func (c *Consistent) GetLeast(key string) (string, error) {
	if len(c.hosts) == 0 {
		return "", ErrNoHosts
	}

	h := hash([]byte(key))
	idx := searchSlice(c.sortedSet, h)

	i := idx
	for {
		host := c.hosts[c.sortedSet[i]]
		if c.loadOK(host) {
			return host, nil
		}
		i++
		if i >= len(c.hosts) {
			i = 0
		}
	}
}

// GetLoad get the load of host
func (c *Consistent) GetLoad(host string) uint64 {
	return c.loadMap[host].Load
}

// Inc host load increases by 1
func (c *Consistent) Inc(host string) {
	if _, ok := c.loadMap[host]; !ok {
		return
	}
	c.loadMap[host].Load += 1
	c.totalLoad += 1
}

// Done host load reduced by 1
func (c *Consistent) Done(host string) {
	if _, ok := c.loadMap[host]; !ok {
		return
	}
	c.loadMap[host].Load -= 1
	c.totalLoad -= 1
}

func (c *Consistent) loadOK(host string) bool {
	avg := float64((c.totalLoad + 1) / uint64(len(c.loadMap)))
	if avg == 0 {
		avg = 1
	}
	avg = math.Ceil(avg * loadFactor)
	return float64(c.loadMap[host].Load)+1 <= avg
}
