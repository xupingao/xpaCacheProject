package singleflight

import "sync"

type Group struct{
	mu sync.Mutex
	m map[string]*call
}
type call struct{
	wg sync.WaitGroup
	val interface{}
	err error
}

func (g *Group)Do(key string,fn func()(interface{},error))(interface{},error){
	g.mu.Lock()
	if g.m==nil{
		g.m=make(map[string]*call)
	}
	if c,ok:=g.m[key];ok{
		g.mu.Unlock()
		c.wg.Wait()
		return c.val,c.err
	}
	call:=&call{}
	call.wg.Add(1)
	g.m[key]=call
	g.mu.Unlock()
	call.val,call.err=fn()
	call.wg.Done()

	g.mu.Lock()
	delete(g.m,key)
	g.mu.Unlock()
	return call.val,call.err
}
