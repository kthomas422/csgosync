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
		"../../test/t1.txt",
		"../../test/t2.txt",
		"../../test/t3.txt",
	}

	testHashes = []string{
		"22596363b3de40b06f981fb85d82312e8c0ed511",
		"648a6a6ffffdaa0badb23b8baf90b6168dd16b3a",
		"09fac8dbfd27bd9b4d23a00eb648aa751789536d",
	}
)

func TestLoadFiles(t *testing.T) {
	files, err := loadFiles("../../test")
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
