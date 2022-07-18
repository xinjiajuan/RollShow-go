package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"fmt"
	"net/http"
	"strconv"
)

type HandlerServer struct {
	ServerInfo Config.Server
}

func StartS3HttpServer(config Config.Yaml) {
	var serverList []Config.Server
	for _, list := range config.ServerList {
		serverList = append(serverList, list)
	}
	for _, server := range serverList {
		webserver := http.Server{
			Addr: ":" + strconv.Itoa(server.ListenPort),
		}
		serverhandler := HandlerServer{}
		serverhandler.ServerInfo = server
		http.Handle("/"+server.Bucket, serverhandler)
		go webserver.ListenAndServe()
		println(server.Name + " is Running to :" + strconv.Itoa(server.ListenPort) + "/" + server.Bucket)
		//go webserver.ListenAndServe()
	}
	select {} //程序堵塞
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
