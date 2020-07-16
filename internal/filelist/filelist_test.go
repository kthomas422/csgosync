// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/filelist/filelist_test.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the functions for testing the fileslist module for the csgo sync application.
*/

package filelist

import "testing"

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
)

func TestLoadFiles(t *testing.T) {
	files, err := loadFiles("test")
	if len(files) != len(testFiles) {
		t.Fatal("Length of files doesn't match, got: ", len(files), " wanted: ", len(testFiles))
	}
	if err != nil {
		t.Fatal("failed to get dir files", err)
	}
	for i := 0; i < len(testFiles); i++ {
		if files[i] != testFiles[i] {
			t.Error("file differs got: ", files[i], " wanted: ", testFiles[i])
		}
	}
}

func TestHashFiles(t *testing.T) {
	hashes, err := hashFiles(testFiles)
	if err != nil {
		t.Fatal("couldn't load hashes")
	}
	for i := 0; i < len(hashes); i++ {
		if hashes[i] != testHashes[i] {
			t.Error("[", testFiles[i], "] hash mismatch, got: ", hashes[i], " wanted: ", testHashes[i])
		}
	}
}
