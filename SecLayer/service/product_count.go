package service

import "sync"

type ProductCountMgr struct {
	ProductCount map[int]int
	lock         sync.RWMutex
}

func NewProductCountMgr() (p *ProductCountMgr) {
	p = &ProductCountMgr{
		ProductCount: make(map[int]int, 128),
	}
	return
}

func (p *ProductCountMgr) Count(productId int) (count int) {
	p.lock.RLock()
	defer p.lock.Unlock()

	count = p.ProductCount[productId]
	return
}

func (p *ProductCountMgr) Add(productId, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	cur, ok := p.ProductCount[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}

	p.ProductCount[productId] = cur
}
