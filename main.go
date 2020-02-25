package main
import(
	"./xpacache"
	"flag"
	"log"
	"net/http"
)
var dbMap=map[string]string{
	"key1":"value1",
	"key2":"value2",
	"key3":"value3",
}
func createGroup() *xpacache.Group {
	group:=xpacache.NewGroup("cacheGroup1",0,xpacache.GetterFunc(func(key string)([]byte,error){
			log.Printf("[Simple dataBase]:Search key:%s\n",key)
			if val,ok:=dbMap[key];ok{
				return []byte(val),nil
			}else{
				return []byte{},nil
			}
	}))
	return group
}
func startCacheServer(addr string, addrs []string, xpa *xpacache.Group) {
	p:=xpacache.NewHttpPool(addr)
	p.Set(addrs...)
	xpa.RegisterPeers(p)
	log.Println("xpaache is running at", addr[7:])
	log.Fatal(http.ListenAndServe(addr[7:],p))
}
func startAPIServer(apiAddr string, xpa *xpacache.Group) {
	http.Handle("/api",http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		key:=r.URL.Query().Get("key")
		view,err:=xpa.Get(key)
		if err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
//		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("API server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr,nil))
}
func main(){
	var port int
	var api bool
	flag.IntVar(&port,"port",8001,"XpaCache server port")
	flag.BoolVar(&api,"api",false,"Start api server or not")
	flag.Parse()
	group:=createGroup()
	peerMap:=map[int]string{
		8001:"http://127.0.0.1:8001",
		8002:"http://127.0.0.1:8002",
		8003:"http://127.0.0.1:8003",
	}
	peerAddrs:=[]string{
		"http://127.0.0.1:8001",
		"http://127.0.0.1:8002",
		"http://127.0.0.1:8003",
	}
	if api{
		go startAPIServer("127.0.0.1:8080",group)
	}
	startCacheServer(peerMap[port],peerAddrs,group)

}
