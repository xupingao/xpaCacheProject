package xpacache
import pb "./xpacachepb"
type PeerGetter interface{
	Get(int *pb.Request,out *pb.Response)(error)
}
type PeerPicker interface{
	PickPeer(key string)(PeerGetter,bool)
}
