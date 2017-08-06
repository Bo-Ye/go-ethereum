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
	"os"
	"testing"
)

var (
	key   = []byte("key-abc")
	value = []byte("value-bcd")
)

//test LDBDatabase
func TestLDBDatabase(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	//test Path
	dbPath := db.Path()
	if dbPath != tmpDir {
		t.Errorf("expected %s got %s", tmpDir, dbPath)
	}
	//test Put
	db.Put(key, value)
	ret1, _ := db.db.Has(key, nil)
	if ret1 != true {
		t.Errorf("expected %t got %t", true, ret1)
	}
	//test Get
	ret2, _ := db.Get(key)
	if !bytes.Equal(ret2, value) {
		t.Errorf("expected %x got %x", value, ret2)
	}
	//test Delete
	db.Delete(key)
	ret3, _ := db.db.Has(key, nil)
	if ret3 {
		t.Errorf("expected %t got %t", false, ret3)
	}
}

func TestLDBDatabaseMeter(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	db.Meter("prefix")
}

//////test table
func TestTable(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	table := NewTable(db, "prefix-")
	//test put
	table.Put(key, value)
	ret1, _ := db.db.Has(append([]byte("prefix-"), key...), nil)
	if !ret1 {
		t.Errorf("expected %t got %t", true, ret1)
	}
	//test get
	ret2, _ := table.Get(key)
	if !bytes.Equal(ret2, value) {
		t.Errorf("expected %s got %s", value, ret2)
	}
	//test delete
	table.Delete(key)
	ret3, _ := db.db.Has(append([]byte("prefix-"), key...), nil)
	if ret3 {
		t.Errorf("expected %t got %t", false, ret3)
	}
}
