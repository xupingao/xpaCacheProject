package xpacache

type ByteView struct{
	b []byte
}
func (b ByteView)Len()int{
	return len(b.b)
}
func (b ByteView)String()(string){
	return string(b.b)
}
func (b ByteView)ByteSlice()[]byte{
	return cloneBytes(b.b)
}
func cloneBytes(in []byte)[]byte{
	res:=make([]byte,len(in))
	copy(res,in)
	return res
}