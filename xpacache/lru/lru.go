package lru
import(
	"container/list"
)
type Cache struct{
	maxBytes int64
	nBytes int64
	cache map[string]*list.Element  //attation
	ll *list.List
	onEvicted func(string,Value)
}
type entry struct{
	key string
	val Value
}
type Value interface{
	Len()int
}
func New(maxByte int64,onEvicted func(string,Value))*Cache{
	return &Cache{
		maxBytes:  maxByte,
		nBytes:    0,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		onEvicted: onEvicted,
	}
}

func (c *Cache)Get(key string)(Value,bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).val,true
	}
	return nil,false
}
func (c *Cache)RemoveOldest()(){
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*entry)
		c.nBytes-=int64(kv.val.Len())+int64(len(kv.key))
		delete(c.cache,kv.key)
		if c.onEvicted!=nil{
			c.onEvicted(kv.key,kv.val)
		}
	}

}

func (c *Cache)Add(key string,val Value)(){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		c.nBytes-=int64(ele.Value.(*entry).val.Len())
		ele.Value=&entry{key,val}
		c.nBytes+=int64(val.Len())
		return
	}
	newEle:=c.ll.PushFront(&entry{key,val})
	c.cache[key]=newEle
	c.nBytes+=int64(len(key))+int64(val.Len())
	for c.maxBytes != 0&&c.nBytes>c.maxBytes{
		c.RemoveOldest()
	}
}

func (c *Cache)Len()int{
	return c.ll.Len()
}
