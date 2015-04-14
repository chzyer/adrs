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
	Msg *dns.DNSMessage
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

func (w *Wiki) Lookup(msg *dns.DNSMessage) (b *dns.DNSMessage, ok bool) {
	var cacheBlock *CacheBlock
	key := msg.Key()
	if key != "" {
		w.m.RLock()
		cacheBlock, ok = w.data[key]
		w.m.RUnlock()
		if ok {
			deadline := cacheBlock.Msg.GetDeadline()
			if time.Now().After(deadline) {
				w.m.Lock()
				delete(w.data, key)
				w.m.Unlock()
				logex.Info("remove cache", `"`+key+`"`, deadline)
				return nil, false
			}
			logex.Info("get cache", `"`+key+`"`)
			return cacheBlock.Msg.Copy(w.pool.Get()), true
		}
	}
	return
}

func (w *Wiki) WriteDown(msg *dns.DNSMessage) bool {
	key := msg.Key()
	if key == "" {
		return false
	}

	if len(msg.Resources) == 0 {
		return false
	}

	w.m.Lock()
	w.data[key] = &CacheBlock{msg}
	w.m.Unlock()
	return true
}
