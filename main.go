package main

import (
	"log"
	//"strings"
	"fmt"
	"net/url"
	"net"
	"context"
	"net/http"
	"flag"
)

func main() {
	var keyfile string
	var certfile string
	var port string
	var httpPort string
	var dir string
	flag.StringVar(&keyfile,"keyfile","","key file path")
	flag.StringVar(&certfile,"certfile","","cert file path")
	flag.StringVar(&port,"port","","https server port")
	flag.StringVar(&httpPort,"httpPort","","http server port")
	flag.StringVar(&dir,"dir","","html directory")
	flag.Parse()

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	httpsL,err := net.Listen("tcp",fmt.Sprintf(":%s",port))
	if err != nil{
		panic(err)
	}

	srv := http.NewServeMux()
	srv.Handle(
			"/",
			http.FileServer(http.Dir(dir)),
	)

	go func(){
		defer cancel()
		fmt.Println("server running")
		err := http.ServeTLS(
			httpsL,
			http.HandlerFunc(
				func(res http.ResponseWriter,req *http.Request){
					srv.ServeHTTP(res,req)
				},
			),
			certfile,keyfile,
		)
		if err != nil{
			fmt.Println(err)
		}
	}()

	httpL,err := net.Listen("tcp",fmt.Sprintf(":%s",httpPort))
	if err != nil{
		panic(err)
	}

	go func(ctx context.Context){
		<-ctx.Done()
		httpL.Close()
		httpsL.Close()
	}(ctx)

	fmt.Println("http server running")
	err = http.Serve(
		httpL,
		http.HandlerFunc(
			func(res http.ResponseWriter,req *http.Request){
				redUrl,err := url.Parse(fmt.Sprintf("https://%s:443",req.Host))
				if err != nil{
					panic(err)
				}
				log.Println(redUrl.String())
				http.Redirect(res,req,redUrl.String(),http.StatusSeeOther)
			},
		),
	)
	cancel()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("server quit")

}
