package xpacache

import (
	"fmt"
	"github.com/golang/groupcache/consistenthash"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	pb "./xpacachepb"
)

var basePath="/_xpacache/"
var defaultReplicas=50
type HttpPool struct{
	mu sync.Mutex
	self string
	basePath string
	peers *consistenthash.Map
	httpGetters map[string]*httpGetter
}
func NewHttpPool(self string)*HttpPool{
	return &HttpPool{
		self:     self,
		basePath: basePath,
	}
}

func (httppool *HttpPool)Flog(format string,v ...interface{}){
	log.Printf("[server %s]:%s",httppool.self,fmt.Sprintf(format,v...))
}
func (p *HttpPool)ServeHTTP(w http.ResponseWriter,request *http.Request){
	if !strings.HasPrefix(request.URL.Path,basePath){
		http.Error(w,"Invalid path",http.StatusBadRequest)
		return
	}
	p.Flog("%s %s",request.Method,request.URL.Path)
	strs:=strings.SplitN(request.URL.Path[len(basePath):],"/",2)
	if len(strs)!=2{
		http.Error(w,"Invalid path",http.StatusBadRequest)
		return
	}
	groupName:=strs[0]
	key:=strs[1]
	group:=GetGroup(groupName);
	if group==nil{
			http.Error(w,"non-exist group "+groupName,http.StatusBadRequest)
			return
	}
	val,err:=group.Get(key)
	if err!=nil{
		http.Error(w,"non-exist key "+key,http.StatusBadRequest)
		return
	}
	res,err:=proto.Marshal(&pb.Response{Value:val.ByteSlice()})
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
//	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(res)
}

type httpGetter struct{
	baseUrl string
}

func (h httpGetter)Get(in *pb.Request,out *pb.Response)(error){
	requestUrl:=fmt.Sprintf("%s%s/%s",h.baseUrl,
		in.Group,in.Key)
	res,err:=http.Get(requestUrl)
	if err!=nil{
		return err
	}
	defer res.Body.Close()   //重要
	if res.StatusCode!=http.StatusOK{
		return fmt.Errorf("Bad request")
	}
	bytes,err:=ioutil.ReadAll(res.Body)
	if err!=nil{
		return fmt.Errorf("Bad response")
	}
	if err=proto.Unmarshal(bytes,out);err!=nil{
		return fmt.Errorf("Decoding failed")
	}
	return nil
}
var _PeerGetter=(*httpGetter)(nil)

func (p *HttpPool)Set(peers ...string){
	//并发保护
	p.mu.Lock()
	defer p.mu.Unlock()
	//一致性哈希表
	p.peers=consistenthash.New(defaultReplicas,nil)
	p.peers.Add(peers...)
	//HttpGetters
	p.httpGetters=make(map[string]*httpGetter,len(peers))
	for _,peerUrl:=range(peers){
		p.httpGetters[peerUrl]=&httpGetter{baseUrl:peerUrl+p.basePath}
	}
}

func (p *HttpPool)PickPeer(key string)(PeerGetter,bool){
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.peers==nil{
		return nil,false
	}
	if peer,ok:=p.httpGetters[p.peers.Get(key)];ok&&peer.baseUrl!=""&&peer.baseUrl!=p.self+p.basePath{
		return peer,true
	}
	return nil,false
}
var _httpPool=(*HttpPool)(nil)
