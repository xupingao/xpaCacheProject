package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte)uint32
type Map struct{
	hash Hash
	replicas int
	keys []int
	hashMap map[int]string
}
func New(replicas int,fn Hash)*Map{
	m:= &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash==nil{
		m.hash=crc32.ChecksumIEEE
	}
	return m
}
func (m *Map) Add(keys ...string) {
	for _,key:=range(keys){
		for i:=0;i<m.replicas;i++{
			hash:=m.hash([]byte(strconv.Itoa(i)+key))
			m.keys= append(m.keys, int(hash))
			m.hashMap[int(hash)]=key
		}
	}
	sort.Ints(m.keys)
}
func (m *Map) Get(key string) string {
	if len(m.keys)==0{
		return ""
	}
	hash:=int(m.hash([]byte(key)))
	idx:=sort.Search(len(m.keys) ,func(i int)bool{
		return m.keys[i]>=hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]  //重要 如果idx==len(keys)说明哈希结果比所有值都大，此时选择m.int[0]

}