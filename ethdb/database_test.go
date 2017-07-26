// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethdb

import (
	"bytes"
	"io/ioutil"
	"testing"
	"fmt"
	"path/filepath"
	"os"
)

var(
	key = []byte("15da97c42b7ed2e1c0c8dab6a6d7e3d9dc0a75580bbc4f1f29c33996d1415dcc")
	value = []byte("Hello world")
)

func TestPath(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "ldbtesttmpfile")
	fmt.Println(tmpDir)
	db, _ := NewLDBDatabase(file, 0, 0)
	dbPath := db.Path()
	exp := file
	fmt.Println(file)
	if dbPath != exp {
		t.Errorf("expected %x got %x", exp, dbPath)
	}
}

func TestPut(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "ldbtesttmpfile")
	db, _ := NewLDBDatabase(file, 0, 0)
	db.Put(key, value)
	ret, _ := db.db.Has(key, nil)
	if ret != true {
		t.Errorf("expected %x got %x", true, ret)
	}
}

func TestGet(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "ldbtesttmpfile")
	db, _ := NewLDBDatabase(file, 0, 0)
	db.db.Put(key, value, nil)
	ret, _ := db.Get(key)
	if !bytes.Equal(ret, value) {
		t.Errorf("expected %x got %x", value, ret)
	}
}

func TestDelete(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "ldbtesttmpfile")
	db, _ := NewLDBDatabase(file, 0, 0)
	db.db.Put(key, value, nil)
	db.Delete(key)
	ret, _ := db.db.Has(key, nil)
	if ret != false {
		t.Errorf("expected %x got %x", true, ret)
	}
}
