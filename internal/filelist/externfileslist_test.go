// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/filelist/externfilelist_test.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file contains the functions for externally testing the fileslist module for the csgo sync application.
*/

package filelist_test

import (
	"testing"

	"github.com/kthomas422/csgosync/internal/filelist"
)

var (
	testFiles = []string{
		"../../test/t1.txt",
		"../../test/t2.txt",
		"../../test/t3.txt",
	}

	testHashes = []string{
		"22596363b3de40b06f981fb85d82312e8c0ed511",
		"648a6a6ffffdaa0badb23b8baf90b6168dd16b3a",
		"09fac8dbfd27bd9b4d23a00eb648aa751789536d",
	}
	testMap = map[string]string{
		"t1.txt": "22596363b3de40b06f981fb85d82312e8c0ed511",
		"t2.txt": "648a6a6ffffdaa0badb23b8baf90b6168dd16b3a",
		"t3.txt": "09fac8dbfd27bd9b4d23a00eb648aa751789536d",
	}
	clientMap = map[string]string{
		"t1.txt": "22596363b3de40b06f981fb85d82312e8c0ed511",
		"t2.txt": "648a6a6ffffdaa0badb23b8baf90b6168dd16b3a",
	}
	deltaMap = []string{
		"t3.txt",
	}
)

func TestGenerateMap(t *testing.T) {
	laMap, err := filelist.GenerateMap("test")
	if err != nil {
		t.Fatal("couldn't get hash map")
	}
	for k, v := range laMap {
		if v != testMap[k] {
			t.Error("hash not found")
		}
	}
}

func TestComparemaps(t *testing.T) {
	delta := filelist.CompareMaps(testMap, clientMap)
	if len(delta) != len(deltaMap) {
		t.Fatal("delta map wrong")
	}
	for i := range deltaMap {
		if delta[i] != deltaMap[i] {
			t.Error("filename mismatch")
		}
	}
}
