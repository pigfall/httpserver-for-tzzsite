package main

import (
	"fmt"
	"net/http"
	"flag"
)

func main() {
	var keyfile string
	var certfile string
	var port string
	flag.StringVar(&keyfile,"keyfile","","key file path")
	flag.StringVar(&certfile,"certfile","","cert file path")
	flag.StringVar(&port,"port","","server port")
	flag.Parse()
	h :=http.FileServer(http.Dir("."))

	fmt.Println("server running")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%s",port),certfile,keyfile,h)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("server quit")
}
