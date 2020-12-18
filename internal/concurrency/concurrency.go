// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/concurrency/concurrency.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file contains concurrency overhead structs and methods for the csgo sync application.
*/

package concurrency

import "sync"

// Token and Semaphore is how the number of threads hashing the files and making http requests will be limited later
type Token struct{}
type Semaphore chan Token

// Wrapper for the concurrency "primitives"
type OverHead struct {
	Wg      *sync.WaitGroup // waitgroup is for keeping track of threads spawned
	HttpSem Semaphore       // HttpSem is a semaphore for limiting the number of http requests sent at once
	FileSem Semaphore       // FileSem is a sempaphore for limiting the number of files opened at once
}

// InitOH returns the initialized concurrency wrapper
func InitOH(maxHttp, maxFile int) *OverHead {
	return &OverHead{
		Wg:      new(sync.WaitGroup),
		FileSem: make(Semaphore, maxFile),
		HttpSem: make(Semaphore, maxHttp),
	}
}
