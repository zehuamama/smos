// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smos

import (
	"hash/crc32"
	"sort"
)

func searchSlice(sortedSet []uint32, key uint32) int {
	idx := sort.Search(len(sortedSet), func(i int) bool {
		return sortedSet[i] >= key
	})

	if idx >= len(sortedSet) {
		idx = 0
	}
	return idx
}

func deleteSlice(sortedSet *[]uint32, val uint32) {
	idx := -1
	l := 0
	r := len(*sortedSet) - 1
	for l <= r {
		m := (l + r) / 2
		if (*sortedSet)[m] == val {
			idx = m
			break
		} else if (*sortedSet)[m] < val {
			l = m + 1
		} else if (*sortedSet)[m] > val {
			r = m - 1
		}
	}
	if idx != -1 {
		*sortedSet = append((*sortedSet)[:idx], (*sortedSet)[idx+1:]...)
	}
}

func hash(key []byte) uint32 {
	return crc32.ChecksumIEEE(key)
}
