package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type HandlerServer struct {
	ServerInfo Config.Server
}

func MakeS3HttpServer(config Config.Yaml) {
	var serverList []Config.Server
	var serverObjectList []*http.Server
	for _, list := range config.ServerList {
		serverList = append(serverList, list)
	}
	for _, server := range serverList {
		serverhandler := HandlerServer{}
		serverhandler.ServerInfo = server
		webserver := http.Server{
			Addr:    ":" + strconv.Itoa(server.ListenPort),
			Handler: serverhandler,
		}
		serverObjectList = append(serverObjectList, &webserver)
		//http.Handle("/"+server.Bucket, serverhandler)
		//go webserver.ListenAndServe()
	}
	RunHttpServer(context.Background(), serverList, serverObjectList)
}

func RunHttpServer(ctx context.Context, serverlist []Config.Server, httpSrv []*http.Server) {
	for i, serverObject := range httpSrv {
		println(serverlist[i].Name + " is Running to :" + strconv.Itoa(serverlist[i].ListenPort) + "/" + serverlist[i].Bucket)
		go serverObject.ListenAndServe()
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sigs:
		for i, serverObject := range httpSrv {
			fmt.Println("Shutting down " + serverlist[i].Name + " instance gracefully...")
			serverObject.Shutdown(ctx)
			fmt.Println("Instance " + serverlist[i].Name + " has exited safely!")
		}
	} //程序堵塞
}

func (webserver HandlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s3ObjectClient, er := MakeClient(webserver.ServerInfo)
	if er != nil {
		fmt.Println(er.Error())
		return
	}
	objectList, er := GetObject(s3ObjectClient, webserver.ServerInfo.Bucket, "", true, false, false)
	if er != nil {
		fmt.Println(er.Error())
		return
	}
	for i, object := range objectList.Info {
		fmt.Fprintf(w, "%d - ", i)
		fmt.Fprintf(w, object.Key)
		fmt.Fprintf(w, " - %d ", object.Size)
		fmt.Fprintln(w, "")
	}
}
