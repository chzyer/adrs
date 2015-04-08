package utils

import "sync"

const POOLSIZE = 512

type BlockPooler interface {
	Get() []byte
	Put([]byte)
}

type BlockPool struct {
	pool sync.Pool
}

func NewBlockPool() *BlockPool {
	p := new(BlockPool)
	p.pool = sync.Pool{New: p.newBlock}
	return p
}
func (p *BlockPool) newBlock() interface{} {
	return make([]byte, POOLSIZE)
}
func (p *BlockPool) Get() []byte {
	return p.pool.Get().([]byte)
}

func (p *BlockPool) Put(b []byte) {
	p.pool.Put(b)
}
