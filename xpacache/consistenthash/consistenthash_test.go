package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistenthash(t *testing.T){
	m:=New(3,Hash(func(in []byte)uint32{
		i,_:=strconv.Atoi(string(in))
		return uint32(i)
	}))
	keys:=[]string{"2","4","6"}
	//2:[2 12 22]
	//4:[4 14 24]
	//6:[6 16 26]
	m.Add(keys...)

	testCases:=map[string]string{
		"3":"4",
		"10":"2",
		"20:":"2",
		"25":"6",
		"5":"6",
		"27":"2",
	}
	for k,v:=range(testCases){
		if m.Get(k)!=v{
			t.Fatalf("Error!The value of key  [%s] should be [%s] ,result is [%s]",k,v,m.Get(k))
		}
	}
	//8:[8 18 28]
	m.Add("8")
	if m.Get("27")!="8"{
		t.Fatalf("Error!The value of key  [%s] should be [%s]","27","8")
	}
}
