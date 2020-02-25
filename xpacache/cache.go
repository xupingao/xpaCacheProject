package xpacache

import (
	"./lru"
	"sync"
)

type cache struct{
	m sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}
func (c *cache)Get(key string)(ByteView,bool) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru == nil {
		return ByteView{nil}, false
	}
	if val, ok := c.lru.Get(key); ok {
		return val.(ByteView), true
	} else {
		return ByteView{nil}, false
	}
}
func (c *cache)Add(key string,val ByteView){
	c.m.Lock()
	defer c.m.Unlock()
	if c.lru==nil{
		c.lru=lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key,val)
}
