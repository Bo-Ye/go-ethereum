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
	"github.com/ethereum/go-ethereum/metrics"
	"io/ioutil"
	"os"
	"testing"
)

var (
	key1       = []byte("key-aaaaa")
	value1     = []byte("value-aaaaa")
	key2       = []byte("key-bbbbb")
	value2     = []byte("value-bbbbb")
	key3       = []byte("key-ccccc")
	value3     = []byte("value-ccccc")
	unknownKey = []byte("unknow")
)

//test LDBDatabase
func TestLDBDatabase(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	defer db.Close()
	//test Path
	dbPath := db.Path()
	if dbPath != tmpDir {
		t.Errorf("LDBDatabase Path: expected %s, got %s", tmpDir, dbPath)
	}
	//test Put
	db.Put(key1, value1)
	ret1, _ := db.db.Has(key1, nil)
	if !ret1{
		t.Error("LDBDatabase Put: expected true, got false")
	}
	//test Get
	ret2, _ := db.Get(key1)
	if !bytes.Equal(ret2, value1) {
		t.Errorf("LDBDatabase Get: expected %s, got %s", value1, ret2)
	}
	//test Delete
	db.Delete(key1)
	ret3, _ := db.db.Has(key1, nil)
	if ret3 {
		t.Error("LDBDatabase Delete: expected false, got true")
	}
}

func TestLDBDatabaseMeter(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	defer db.Close()
	//enable metrics
	metrics.Enabled = true
	db.Meter("prefix")
	//test putTimer && writeMeter
	db.Put(key1, value1)
	db.Put(key2, value2)
	ret := db.putTimer.Count()
	if ret != 2 {
		t.Errorf("putTimer: expected %d, got %d", 2, ret)
	}
	min := db.putTimer.Min()
	if min == 0 {
		t.Error("putTimer: expected min time larger than zero, got 0")
	}
	mean := db.putTimer.Mean()
	if mean == 0 {
		t.Error("putTimer: expected mean time larger than zero, got 0")
	}
	max := db.putTimer.Max()
	if max == 0 {
		t.Error("putTimer: expected max time larger than zero, got 0")
	}
	totalBytes := int64(len(value1) + len(value2))
	ret = db.writeMeter.Count()
	if totalBytes != ret {
		t.Errorf("writeMeter: expected %d, got %d", totalBytes, ret)
	}
	//test getTimer & missMeter & readMeter
	db.Get(key1)
	db.Get(unknownKey)
	ret = db.getTimer.Count()
	if ret != 2 {
		t.Errorf("getTimer: expected %d, got %d", 2, ret)
	}
	min = db.getTimer.Min()
	if min == 0 {
		t.Error("getTimer: expected min time larger than zero, got 0")
	}
	mean = db.getTimer.Mean()
	if mean == 0 {
		t.Error("getTimer: expected mean time larger than zero, got 0")
	}
	max = db.getTimer.Max()
	if max == 0 {
		t.Error("getTimer: expected max time larger than zero, got 0")
	}
	ret = db.missMeter.Count()
	if ret != 1 {
		t.Errorf("missMeter: expected %d got %d", 1, ret)
	}
	totalBytes = int64(len(value1))
	ret = db.readMeter.Count()
	if totalBytes != ret {
		t.Errorf("readMeter: expected %d got %d", totalBytes, ret)
	}
	//test delTimer
	db.Delete(key1)
	db.Delete(key2)
	ret = db.delTimer.Count()
	if ret != 2 {
		t.Errorf("delTimer: expected %d got %d", 2, ret)
	}
	min = db.delTimer.Min()
	if min == 0 {
		t.Error("delTimer: expected min time larger than zero, got 0")
	}
	mean = db.delTimer.Mean()
	if mean == 0 {
		t.Error("delTimer: expected mean time larger than zero, got 0")
	}
	max = db.delTimer.Max()
	if max == 0 {
		t.Error("delTimer: expected max time larger than zero, got 0")
	}
}

////test Batch
func TestBatch(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	defer db.Close()
	//write batch
	batch := db.NewBatch()
	batch.Put(key1, value1)
	batch.Put(key2, value2)
	batch.Put(key3, value3)
	batch.Write()
	ret, _ := db.Get(key1)
	if !bytes.Equal(ret, value1) {
		t.Errorf("batch: expected %x, got %x", value1, ret)
	}
	ret, _ = db.Get(key2)
	if !bytes.Equal(ret, value2) {
		t.Errorf("batch: expected %x, got %x", value2, ret)
	}
	ret, _ = db.Get(key3)
	if !bytes.Equal(ret, value3) {
		t.Errorf("batch: expected %x, got %x", value3, ret)
	}
}

//////test table
func TestTable(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "ldbtesttmpdir")
	defer os.RemoveAll(tmpDir)
	db, _ := NewLDBDatabase(tmpDir, 0, 0)
	table := NewTable(db, "prefix-")
	//test put
	table.Put(key1, value1)
	ret1, _ := db.db.Has(append([]byte("prefix-"), key1...), nil)
	if !ret1 {
		t.Errorf("expected %t, got %t", true, ret1)
	}
	//test get
	ret2, _ := table.Get(key1)
	if !bytes.Equal(ret2, value1) {
		t.Errorf("expected %s, got %s", value1, ret2)
	}
	//test delete
	table.Delete(key1)
	ret3, _ := db.db.Has(append([]byte("prefix-"), key1...), nil)
	if ret3 {
		t.Errorf("expected %t, got %t", false, ret3)
	}
	//test tableBatch
	tableBatch := table.NewBatch()
	tableBatch.Put(key1, value1)
	tableBatch.Put(key2, value2)
	tableBatch.Put(key3, value3)
	tableBatch.Write()
	ret, _ := db.db.Has(append([]byte("prefix-"), key1...), nil)
	if !ret {
		t.Error("tableBatch : expected true, got false")
	}
	ret, _ = db.db.Has(append([]byte("prefix-"), key2...), nil)
	if !ret {
		t.Error("tableBatch: expected true, got false")
	}
	ret, _ = db.db.Has(append([]byte("prefix-"), key3...), nil)
	if !ret {
		t.Error("tableBatch: expected true, got false")
	}
}
