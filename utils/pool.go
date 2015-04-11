package utils

import (
	"sync"
	"unsafe"

	"gopkg.in/logex.v1"
)

const BLOCK_CAP = 512

type Block struct {
	All      []byte
	Length   int
	pool     *BlockPool
	recycled bool
}

func NewBlockWithByte(b []byte) *Block {
	return &Block{
		All:    b,
		Length: len(b),
	}
}

func NewBlock(pool *BlockPool) *Block {
	b := &Block{
		All:      make([]byte, BLOCK_CAP),
		Length:   0,
		pool:     pool,
		recycled: false,
	}
	return b
}

func (b *Block) Init() {
	b.Length = 0
	b.recycled = false
	logex.Debug(unsafe.Pointer(b), "init!")
}

func (b *Block) Write(bytes []byte) (int, error) {
	n := copy(b.All[b.Length:], bytes)
	b.Length += n
	return n, nil
}

func (b *Block) Recycle() {
	if b.recycled {
		return
	}
	if b.pool == nil {
		return
	}

	logex.Debug(unsafe.Pointer(b), "recycled!")
	b.pool.Put(b)
	b.recycled = true
}

func (b *Block) Len() int {
	return b.Length
}

func (b *Block) Bytes() []byte {
	if b.Length > cap(b.All) {
		return b.All
	}
	return b.All[:b.Length]
}

type BlockPooler interface {
	Get() *Block
	Put(*Block)
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
	return NewBlock(p)
}

func (p *BlockPool) Get() *Block {
	block := p.pool.Get().(*Block)
	block.Init()
	return block
}

func (p *BlockPool) Put(b *Block) {
	p.pool.Put(b)
}
