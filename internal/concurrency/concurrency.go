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
