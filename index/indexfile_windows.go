// Copyright 2016 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0.

 //go:build windows

package index

import (
	"fmt"
	"math"
	"os"
)

type fileIndex struct {
	name string
	data []byte
}

func (f *fileIndex) Read(off, sz uint32) ([]byte, error) {
	end := uint64(off) + uint64(sz)
	if end > uint64(len(f.data)) {
		return nil, fmt.Errorf("out of bounds: %d, len %d, name %s", end, len(f.data), f.name)
	}
	return f.data[off:uint32(end)], nil
}

func (f *fileIndex) Name() string {
	return f.name
}

func (f *fileIndex) Size() (uint32, error) {
	return uint32(len(f.data)), nil
}

func (f *fileIndex) Close() {}

func NewIndexFile(f *os.File) (IndexFile, error) {
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	sz := fi.Size()
	if sz < 0 || sz >= math.MaxUint32 {
		return nil, fmt.Errorf("file %s too large: %d", f.Name(), sz)
	}

	data := make([]byte, int(sz))
	if _, err := f.ReadAt(data, 0); err != nil && len(data) > 0 {
		return nil, err
	}

	return &fileIndex{
		name: f.Name(),
		data: data,
	}, nil
}