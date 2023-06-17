// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package muxpatterns

// A mapping is a set of key-value pairs.
// An zero mapping is empty and ready to use.
//
// Mappings try to pick a representation that makes [mapping.find] most efficient.
type mapping[K comparable, V any] struct {
	s []entry[K, V] // for few pairs
	m map[K]V       // for many pairs
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

// maxSlice is the maximum number of pairs for which a slice is used.
// It is a variable for benchmarking.
var maxSlice int = 8

// add adds a key-value pair to the mapping.
func (h *mapping[K, V]) add(k K, v V) {
	if h.m == nil && len(h.s) < maxSlice {
		h.s = append(h.s, entry[K, V]{k, v})
	} else {
		if h.m == nil {
			h.m = map[K]V{}
			for _, e := range h.s {
				h.m[e.key] = e.value
			}
			h.s = nil
		}
		h.m[k] = v
	}
}

// find returns the value corresponding to the given key.
// The second return value is false if there is no value
// with that key.
func (h *mapping[K, V]) find(k K) (v V, found bool) {
	if h == nil {
		return v, false
	}
	if h.m != nil {
		v, found = h.m[k]
		return v, found
	}
	for _, e := range h.s {
		if e.key == k {
			return e.value, true
		}
	}
	return v, false
}

// pairs calls f for every pair in the mapping.
// If f, returns a non-nil error, pairs returns immediately with the same error.
func (h *mapping[K, V]) pairs(f func(k K, v V) error) error {
	if h == nil {
		return nil
	}
	if h.m != nil {
		for k, v := range h.m {
			if err := f(k, v); err != nil {
				return err
			}
		}
	} else {
		for _, e := range h.s {
			if err := f(e.key, e.value); err != nil {
				return err
			}
		}
	}
	return nil
}
