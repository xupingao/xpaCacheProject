package xpacache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var dbMap=map[string]string{
	"key1":"value1",
	"key2":"value2",
	"key3":"value3",
}
func TestGroup(t *testing.T){
	db:=func(key string)([]byte,error){
		log.Printf("[db]SearchKey %s",key)
		if val,ok:=dbMap[key];ok{
			return []byte(val),nil
		}
		return []byte{},fmt.Errorf("%s is not exist",key)
	}
	group:=NewGroup("testGroup",0,GetterFunc(db))
	for k,v:=range(dbMap){
		for i:=0;i<3;i++{
			if val,err:=group.Get(k);err!=nil||!reflect.DeepEqual([]byte(v),val.b){
				t.Fatalf("get wrong value!")
			}
		}
	}
	if _,err:=group.Get("key4");err==nil{
		t.Fatalf("get none-exist value")
	}
}
func TestGetter(t *testing.T){
	f:=func(key string)([]byte,error){
		return []byte("helloworld"),nil
	}
	var getterFunc Getter=GetterFunc(f)
	//res,err:=getterFunc.Get("test")
	//if err!=nil||string(res)!="helloworld"{
	//	t.Fatalf("callbackFailed")
	//}
	expect:=[]byte("helloworld")
	if res,_:=getterFunc.Get("test");!reflect.DeepEqual(expect,res){
		t.Fatalf("callbackFailed")
	}
}
