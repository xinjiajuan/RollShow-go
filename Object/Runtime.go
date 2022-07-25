package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"context"
	"fmt"
	"github.com/klarkxy/gohtml"
	"github.com/minio/minio-go/v7"
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
		mux := http.NewServeMux()
		mux.Handle("/"+server.Bucket, serverhandler)
		webserver := http.Server{
			Addr:    ":" + strconv.Itoa(server.ListenPort),
			Handler: mux,
		}
		serverObjectList = append(serverObjectList, &webserver)
		//serverObjectList = append(serverObjectList, )
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
	//程序堵塞
	case <-sigs: //检测Ctrl+c退出程序命令
		for i, serverObject := range httpSrv {
			fmt.Println("Shutting down " + serverlist[i].Name + " instance gracefully...")
			serverObject.Shutdown(ctx) //平滑关闭Http Server线程
			fmt.Println("Instance " + serverlist[i].Name + " has exited safely!")
		}
	}
}

func (webserver HandlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s3ObjectClient, er := MakeClient(webserver.ServerInfo)
	if er != nil {
		fmt.Println(er.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	objectList := s3ObjectClient.ListObjects(ctx, webserver.ServerInfo.Bucket, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	var i = 0 //对象计数
	var InfoList []Config.ObjectInfo
	for object := range objectList {
		if object.Err != nil {
			fmt.Println(er.Error())
			return
		}
		i++
		Info := Config.ObjectInfo{}
		Info.Num = i
		Info.Key = object.Key
		Info.Size = object.Size
		Info.ETag = object.ETag
		Info.LsatModified = object.LastModified
		InfoList = append(InfoList, Info)
	}
	fmt.Fprintf(w, makeHomePageHtml(InfoList, webserver.ServerInfo))
}

func makeHomePageHtml(infolist []Config.ObjectInfo, serverInfo Config.Server) string {
	//bootstrap := bootstrap.Bootstrap()
	//bootstrap.Body().H2().Text(serverInfo.Name + " - Bucket: " + serverInfo.Bucket)
	//bootstrap.Body().Hr()
	//for _, object := range infolist {
	//	bootstrap.Body().Tag("a").Href("d/" + object.Key).Text(object.Key).Target("_black").Br()
	//}
	////fmt.Println(bootstrap.String())
	//return bootstrap.String()
	//html构建
	bootstrap := gohtml.NewHtml()
	bootstrap.Html().Lang("zh-CN")
	// Meta部分
	bootstrap.Meta().Charset("utf-8")
	bootstrap.Meta().Http_equiv("X-UA-Compatible").Content("IE=edge")
	bootstrap.Meta().Name("viewport").Content("width=device-width, initial-scale=1")
	// Css引入
	bootstrap.Link().Href("https://cdnjs.cloudflare.com/ajax/libs/bootstrap/5.2.0/css/bootstrap.min.css").Rel("stylesheet")
	// Js引入
	bootstrap.Script().Src("https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0/jquery.min.js")
	bootstrap.Script().Src("https://cdnjs.cloudflare.com/ajax/libs/bootstrap/5.2.0/js/bootstrap.min.js")
	// Head
	bootstrap.Head().Title().Text("S3 Server " + serverInfo.Name + " - " + serverInfo.Bucket)
	// Body
	divframe := bootstrap.Body().Div()
	divframe.Class("container-md")
	divframe.H1().Text("S3 Object Storage File WEB Browser").Class("container text-center")
	divframe.Body().Hr()
	divframe.Body().H3().Small().Class("text-muted").Text("The current bucket is " + serverInfo.Bucket).Hr()
	tablediv := divframe.Body().Div().Class("text-center")
	for _, object := range infolist {
		rowtable := tablediv.Body().Div().Class("row")
		rowtable.Body().Div().Class("col-1").Span().Class("badge bg-secondary").Text(strconv.Itoa(object.Num))
		rowtable.Body().Div().Class("col-9").Align("left").Tag("a").Href("d/" + object.Key).Text(object.Key).Target("_black")
		rowtable.Body().Div().Class("col-2").Text(getObjectSizeSuitableUnit(object.Size))
		//divframe.Body()
	}
	//println(getObjectSizeSuitableUnit(165))
	return bootstrap.String()
}
