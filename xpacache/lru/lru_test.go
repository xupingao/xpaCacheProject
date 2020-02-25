package lru
import(
	"reflect"
	"testing"
)
type String string

func (s String)Len()int{
	return len(string(s))
}
func TestGet(t *testing.T){
	lru:=New(int64(0),nil)
	lru.Add("key1",String("1234"))
	if v,ok:=lru.Get("key1");!ok||string(v.(String))!="1234"{
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _,ok:=lru.Get("key2");ok {
		t.Fatalf("cache miss key2 falied")
	}
}
func TestRemoveoldest(t *testing.T){
	k1,k2,k3:="k1","k2","k3"
	v1,v2,v3:="v1","v2","v3"
	cap:=len(k1+k2+v1+v2)
	lru:=New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))
	if _,ok:=lru.Get("k1");ok{
		t.Fatalf("RemoveOldest key1 failed")
	}
}
func TestOnEvicted(t *testing.T){
	keys:=make([]string,0)
	callback:=func(key string,value Value)(){
		keys= append(keys, key)
	}
	lru:=New(int64(10),callback)
	lru.Add("key1",String("1234"))
	lru.Add("k2",String("k2"))
	lru.Add("k3",String("k3"))
	lru.Add("k4",String("k4"))

	expect:=[]string{"key1","k2"}
	if !reflect.DeepEqual(expect,keys){
		t.Fatalf("Call OnEvicted falied")
	}
}
func TestAdd(t *testing.T){
	lru:=New(int64(0),nil)
	lru.Add("k1",String("k1"))
	lru.Add("k2",String("k2"))
	if lru.nBytes!=int64(len("k2")+len("k2")){
		t.Fatal("expect 4 but got ",lru.nBytes)
	}
}