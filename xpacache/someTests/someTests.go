package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func main(){

	http.ListenAndServe("localhost:8001",http.HandlerFunc(func(w http.ResponseWriter,req *http.Request)(){
		fmt.Fprintf(w,"fmt.Printf")
		io.WriteString(w,"io.Write")
		w.Write([]byte("w.write"))
		bufW:=bufio.NewWriter(w)
		fmt.Fprint(bufW,"buffio1")
		bufW.Write([]byte("bufio2"))
		bufW.Flush()
	}))
}
