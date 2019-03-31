package service

import "sync"

type UserBuyHistory struct {
	history map[int]int
	Lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductByCount(UserId int) int {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	count, _ := p.history[UserId]
	return count
}

func (p *UserBuyHistory) Add(UserId int, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.history[UserId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.history[UserId] = cur
}
