package service

import "sync"

type ProductCountMgr struct {
	ProductCount map[int64]int
	lock         sync.RWMutex
}

func NewProductCountMgr() (p *ProductCountMgr) {
	p = &ProductCountMgr{
		ProductCount: make(map[int64]int, 128),
	}
	return
}

func (p *ProductCountMgr) Count(Id int64) (count int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	count = p.ProductCount[Id]
	return
}

func (p *ProductCountMgr) Add(Id int64, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	cur, ok := p.ProductCount[Id]
	if !ok {
		cur = count
	} else {
		cur += count
	}

	p.ProductCount[Id] = cur
}
