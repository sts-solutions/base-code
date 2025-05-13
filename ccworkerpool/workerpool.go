package ccworkerpool

import "sync"

type Pool struct {
	inProgress chan struct{}
	wg         *sync.WaitGroup
	max        int
}

func New(max int) Pool {
	return Pool{
		inProgress: make(chan struct{}, max),
		wg:         &sync.WaitGroup{},
		max:        max,
	}
}

func (m Pool) Add() {
	m.inProgress <- struct{}{}
	m.wg.Add(1)
}

func (m Pool) Remove() {
	<-m.inProgress
	m.wg.Done()
}

func (m Pool) InProgress() int {
	return len(m.inProgress)
}

func (m Pool) Avaliable() int {
	return m.max - m.InProgress()
}

func (m Pool) WaitToFinish() {
	m.wg.Wait()
}
