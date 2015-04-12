package wiki

import (
	"sync"
	"time"

	"github.com/chzyer/adrs/conf"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type CacheBlock struct {
	Block    *utils.Block
	Deadline time.Time
}

type Wiki struct {
	TTLRatio time.Duration
	data     map[string]*CacheBlock
	m        sync.RWMutex
	pool     *utils.BlockPool
}

func NewWiki(c *conf.Config, pool *utils.BlockPool) *Wiki {
	return &Wiki{
		TTLRatio: time.Duration(c.TTLRatio),
		data:     make(map[string]*CacheBlock, 1024),
		pool:     pool,
	}
}

func (w *Wiki) Lookup(msg *dns.DNSMessage) (b *utils.Block, ok bool) {
	var cacheBlock *CacheBlock
	key := msg.Key()
	if key != "" {
		w.m.RLock()
		cacheBlock, ok = w.data[key]
		w.m.RUnlock()
		if ok {
			if time.Now().After(cacheBlock.Deadline) {
				w.m.Lock()
				delete(w.data, key)
				w.m.Unlock()
				logex.Info("remove cache", `"`+key+`"`, cacheBlock.Deadline)
				return nil, false
			}
			logex.Info("get cache", `"`+key+`"`, "ttl:", int(cacheBlock.Deadline.Sub(time.Now()).Seconds()))
			b = w.pool.Get()
			cacheBlock.Block.CopyTo(b)
			return b, true
		}
	}
	return
}

func (w *Wiki) WriteDown(msg *dns.DNSMessage, block *utils.Block) bool {
	key := msg.Key()
	if key == "" {
		return false
	}
	b := w.pool.Get()
	block.CopyTo(b)

	if len(msg.Resources) == 0 {
		return false
	}

	deadline := time.Now().Add(time.Duration(msg.Resources[0].TTL) / w.TTLRatio * time.Second)

	w.m.Lock()
	w.data[key] = &CacheBlock{b, deadline}
	w.m.Unlock()
	return true
}
