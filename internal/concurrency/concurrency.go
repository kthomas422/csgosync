// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/concurrency/concurrency.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains concurrency overhead stucts and methods for the csgo sync application.
*/

package concurrency

import "sync"

type Token struct{}
type Semaphore chan Token

type OverHead struct {
	Wg      *sync.WaitGroup
	HttpSem Semaphore
	FileSem Semaphore
}

func InitOH(maxHttp, maxFile int) *OverHead {
	oh := new(OverHead)
	oh.Wg = new(sync.WaitGroup)
	oh.FileSem = make(Semaphore, maxFile)
	oh.HttpSem = make(Semaphore, maxHttp)
	return oh
}
