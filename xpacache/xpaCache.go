package xpacache

import (
	"fmt"
	"log"
	"sync"
	"./singleflight"
	pb "./xpacachepb"
)

type Getter interface{
	Get(string)([]byte,error)
}
type GetterFunc func(string)([]byte,error)

func (f GetterFunc)Get(key string)([]byte,error){
	return f(key)
}

type Group struct{
	name string
	getter Getter
	mainCache cache
	peers PeerPicker
	loader *singleflight.Group
}
var(
	mu sync.RWMutex
	Groups=make(map[string]*Group)
)
func NewGroup(name string,cacheBytes int64,getter Getter)*Group{
	if getter==nil{
		panic("nil getter")
	}
	if _,ok:=Groups[name];ok{
		panic("name already exists")
	}
	mu.Lock()
	defer mu.Unlock()
	group:=&Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes:cacheBytes},
		loader: &singleflight.Group{},
	}
	Groups[name]=group
	return group
}
func (g *Group)RegisterPeers(p *HttpPool){
	if g.peers!=nil{
		panic("Gegister function called more then once")
	}
	g.peers=p
}
func GetGroup(name string)*Group{
	mu.RLock()
	defer mu.RUnlock()
	return Groups[name]
}

func (g *Group)Get(key string)(ByteView,error){
	if key==""{
		return ByteView{},fmt.Errorf("Invalid key")
	}
	if val,ok:=g.mainCache.Get(key);ok{
		log.Printf("[xpaCache]:Get key \"%s\" directly from cache",key)
		return val,nil
	}
	return g.load(key)
}
func (g *Group)load(key string)(ByteView,error){
	loaderVal,err:=g.loader.Do(key,func()(interface{},error){
		if g.peers!=nil{
			if peer,ok:=g.peers.PickPeer(key);ok{
				if res,err:=g.getFromPeer(key,peer);err==nil{
					log.Printf("[xpaCache]:Get key \"%s\" from peer",key)
					return res,nil
				}else{
					log.Println("[xpaCache]: Failed to get from peer:", err)
				}
			}
		}
		return g.getLocally(key)
	})
	if err!=nil{
		return ByteView{},err
	}
	return loaderVal.(ByteView),nil
}
func (g *Group)getFromPeer(key string,picker PeerGetter)(ByteView,error){
	req:=&pb.Request{
		Key:                 key,
		Group:               g.name,
	}
	res:=new(pb.Response)
	err:=picker.Get(req,res)
	if err!=nil{
		return ByteView{},err
	}
	return ByteView{res.Value},nil
}
func (g *Group)getLocally(key string)(ByteView,error){
	val,err:=g.getter.Get(key)
	if err!=nil{
		return ByteView{},err
	}
	//for safe
	value:=ByteView{cloneBytes(val)}
	log.Printf("[xpaCache]:Get key \"%s\" from local database",key)
	g.populateCache(key,value)
	return value,nil
}
func (g *Group)populateCache(key string,val ByteView){
	g.mainCache.Add(key,val)
}