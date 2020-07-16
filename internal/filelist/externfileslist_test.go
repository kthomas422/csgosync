// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/filelist/externfilelist_test.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the functions for externally testing the fileslist module for the csgo sync application.
*/

package filelist_test

import (
	"testing"

	"github.com/kthomas422/csgosync/internal/filelist"
)

var (
	testFiles = []string{
		"test/t1.txt",
		"test/t2.txt",
		"test/t3.txt",
	}

	testHashes = []string{
		"a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447",
		"d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26",
		"d9014c4624844aa5bac314773d6b689ad467fa4e1d1a50a1b8a99d5a95f72ff5",
	}
	testMap = map[string]string{
		"test/t1.txt": "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447",
		"test/t2.txt": "d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26",
		"test/t3.txt": "d9014c4624844aa5bac314773d6b689ad467fa4e1d1a50a1b8a99d5a95f72ff5",
	}
	clientMap = map[string]string{
		"test/t1.txt": "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447",
		"test/t2.txt": "d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26",
	}
	deltaMap = []string{
		"test/t3.txt",
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
	for i, _ := range deltaMap {
		if delta[i] != deltaMap[i] {
			t.Error("filename mismatch")
		}
	}
}
