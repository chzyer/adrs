package utils

import (
	"sync"
	"unsafe"

	"gopkg.in/logex.v1"
)

const BLOCK_CAP = 512

type Block struct {
	start    []byte
	Block    []byte
	Length   int
	pool     *BlockPool
	recycled bool
}

func NewBlockWithByte(b []byte) *Block {
	return &Block{
		Block:  b,
		Length: len(b),
	}
}

func NewBlock(pool *BlockPool) *Block {
	// tcp header
	block := make([]byte, BLOCK_CAP+2)
	b := &Block{
		start:    block,
		Block:    block[2:],
		Length:   0,
		pool:     pool,
		recycled: false,
	}
	return b
}

// for write
func (b *Block) PlusHeaderBlock() []byte {
	return b.start
}

func (b *Block) SetLengthPlusHeader(l int) {
	b.Length = l - 2
}

// for read
func (b *Block) PlusHeaderBytes() []byte {
	copy(b.start[:2], Uint16To(uint16(b.Length)))
	return b.start[:b.PlusHeaderLength()]
}

func (b *Block) PlusHeaderLength() int {
	return b.Length + 2
}

func (b *Block) Init() {
	b.Length = 0
	b.recycled = false
	logex.Debug(unsafe.Pointer(b), "init!")
}

func (b *Block) Write(bytes []byte) (int, error) {
	n := copy(b.Block[b.Length:], bytes)
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
	if b.Length > cap(b.Block) {
		return b.Block
	}
	return b.Block[:b.Length]
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
