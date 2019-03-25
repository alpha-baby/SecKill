package service

import "sync"

type UserBuyHistory struct {
	history map[int]int
	Lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductByCount(productId int) int {
	p.Lock.RLock()
	defer p.Lock.Unlock()

	count, _ := p.history[productId]
	return count
}

func (p *UserBuyHistory) Add(productId, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.history[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.history[productId] = cur
}
