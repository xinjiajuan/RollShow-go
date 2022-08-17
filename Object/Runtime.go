package Object

import (
	"context"
	"fmt"
	"github.com/klarkxy/gohtml"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"rollshow/Object/Config"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type HandlerServer struct {
	ServerInfo Config.Server
}

//生成http服务对象
func MakeS3HttpServer(config Config.Yaml) {
	var serverList []Config.Server
	var serverObjectList []*http.Server
	for _, list := range config.ServerList {
		if list.Enable {
			serverList = append(serverList, list)
		}
	}
	for _, server := range serverList {
		serverhandler := HandlerServer{}
		serverhandler.ServerInfo = server
		webserver := http.Server{
			Addr:    ":" + strconv.Itoa(server.ListenPort),
			Handler: serverhandler,
		}
		serverObjectList = append(serverObjectList, &webserver)
	}
	RunHttpServer(serverList, serverObjectList)
}

//运行http服务
func RunHttpServer(serverlist []Config.Server, httpSrv []*http.Server) {
	for i, serverObject := range httpSrv {
		println(serverlist[i].Name + " is Running to :" + strconv.Itoa(serverlist[i].ListenPort) + "/" + serverlist[i].Bucket)
		go serverObject.ListenAndServe() //协程并发监听http服务
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	//程序堵塞
	case <-sigs: //检测Ctrl+c退出程序命令
		for i, serverObject := range httpSrv {
			fmt.Println("Shutting down " + serverlist[i].Name + " instance gracefully...")
			serverObject.Shutdown(context.Background()) //平滑关闭Http Server线程
			fmt.Println("Instance " + serverlist[i].Name + " has exited safely!")
		}
	}
}

//处理请求
func (webserver HandlerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//获取url和切分
	urlArray := strings.Split(r.URL.String(), "/")
	//判断url是哪一种请求
	//网站图标
	if r.URL.RequestURI() == "/favicon.ico" {
		w.Header().Set("Location", webserver.ServerInfo.Options.Favicon)
		w.WriteHeader(301)
	}
	//1.根目录，404
	if !strings.EqualFold(urlArray[1], webserver.ServerInfo.Bucket) {
		fmt.Fprintln(w, ErrorPage_404(webserver.ServerInfo.Bucket))
		return
	}
	//2.
	if len(urlArray) > 2 && !strings.EqualFold(urlArray[2], "d") {
		fmt.Fprintln(w, ErrorPage_404(webserver.ServerInfo.Bucket))
		return
	}
	//3.判断/bucket/d/后面有没有东西 404
	if len(urlArray) == 3 && strings.EqualFold(urlArray[2], "d") {
		fmt.Fprintln(w, ErrorPage_404(webserver.ServerInfo.Bucket))
		return
	}

	//文件下载
	if len(urlArray) > 3 && strings.EqualFold(urlArray[2], "d") {
		s3ObjectClient, er := MakeClient(webserver.ServerInfo)
		if er != nil {
			fmt.Println(er.Error())
			fmt.Fprintln(w, er.Error())
			return
		}
		str := strings.SplitN(r.URL.String(), "/", 4)
		enEscapeUrl, _ := url.QueryUnescape(str[3])
		objectStream, er := s3ObjectClient.GetObject(
			context.Background(),
			webserver.ServerInfo.Bucket,
			enEscapeUrl,
			minio.GetObjectOptions{})
		if er != nil {
			fmt.Println(er)
			fmt.Fprintln(w, er.Error())
			return
		}
		// 资源关闭
		defer objectStream.Close()
		info, err := objectStream.Stat()
		if err != nil {
			log.Println("sendFile1", err.Error())
			http.NotFound(w, r)
			return
		}
		w.Header().Add("Accept-ranges", "bytes")
		w.Header().Add("Content-Disposition", "attachment; filename="+urlArray[len(urlArray)-1])
		w.Header().Add("Access-Control-Allow-Origin", webserver.ServerInfo.Options.AccessControlAllowOrigin)
		w.Header().Add("content-type", info.ContentType)
		var start, end int64
		//fmt.Println(request.Header,"\n")
		if ra := r.Header.Get("Range"); ra != "" {
			if strings.Contains(ra, "bytes=") && strings.Contains(ra, "-") {

				fmt.Sscanf(ra, "bytes=%d-%d", &start, &end)
				if end == 0 {
					end = info.Size - 1
				}
				if start > end || start < 0 || end < 0 || end >= info.Size {
					w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
					log.Println("sendFile2 start:", start, "end:", end, "size:", info.Size)
					return
				}
				w.Header().Add("Content-Length", strconv.FormatInt(end-start+1, 10))
				w.Header().Add("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, info.Size))
				w.WriteHeader(http.StatusPartialContent)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			// 非断点续传
			fmt.Println(time.Now().Format(time.UnixDate), r.URL.RequestURI(), r.Proto, r.Host, r.UserAgent(), r.URL.Query().Get("mz_id"))
			println()
			w.Header().Add("Content-Length", strconv.FormatInt(info.Size, 10))
			start = 0
			end = info.Size - 1
		}
		_, err = objectStream.Seek(start, 0)
		// add compare
		if start == (end - start + 1) {
			return
		}
		if err != nil {
			log.Println("sendFile3", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		n := 512
		buf := make([]byte, n)
		for {
			if end-start+1 < int64(n) {
				n = int(end - start + 1)
			}
			//原生 io
			_, er := io.CopyBuffer(w, objectStream, buf)
			if er != nil {
				//log.Println(err, start, end, info.Size(), n)
				return
			}
			start += int64(n)
			if start >= end+1 {
				return
			}
		}

	}
	//显示首页
	html := HomePage(webserver)
	fmt.Fprintln(w, html)

}

//生成404页面
func ErrorPage_404(bucket string) string {
	html := gohtml.NewHtml()
	html.Head().Title().Text("Url Error")
	html.Meta().Charset("utf-8")
	html.Body().H3().Text("404 链接不正确，请添加Bucket路径")
	html.Body().A().Href("/" + bucket).Text(bucket)
	return html.String()
}

//生成主页
func HomePage(webserver HandlerServer) string {
	s3ObjectClient, er := MakeClient(webserver.ServerInfo)
	if er != nil {
		fmt.Println(er.Error())
		return "Create Client:" + er.Error()
	}
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	objectList := s3ObjectClient.ListObjects(ctx, webserver.ServerInfo.Bucket, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	//var i = 0 //对象计数
	var InfoList []Config.ObjectInfo
	for object := range objectList {
		if object.Err != nil {
			fmt.Println(object.Err.Error())
			return "S3 return:" + object.Err.Error()
		}
		//i++
		Info := Config.ObjectInfo{}
		//Info.Num = i
		Info.Key = object.Key
		Info.Size = object.Size
		Info.ETag = object.ETag
		Info.LsatModified = object.LastModified
		InfoList = append(InfoList, Info)
	}
	return makeHomePageHtml(InfoList, webserver.ServerInfo)
}

//生成主页html对象
func makeHomePageHtml(infolist []Config.ObjectInfo, serverInfo Config.Server) string {
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
	divframe := bootstrap.Body().Div()                                                                    //声明容器
	divframe.Class("container-md")                                                                        //容器样式
	divframe.H1().Text("RollShow - S3 Object Server").Class("container text-center")                      //大标题
	divframe.Body().Hr()                                                                                  //分割线
	divframe.Body().Div().Class("alert alert-success").Text("The current bucket is " + serverInfo.Bucket) //桶标识
	//ul := divframe.Body().Ul().Class("list-group")        //文件列表

	table := divframe.Body().Div().Class("rounded border border-success p-2 mb-2").Table().Class("table border-primary table-hover")
	tr := table.Body().Thead().Class("table-light").Tr()
	tr.Body().Th().Attr("scope", "col").Text("#")
	tr.Body().Th().Attr("scope", "col").Text("Object")
	tr.Body().Th().Attr("scope", "col").Text("Info")
	tb := table.Body().Tbody().Class("table-group-divider")
	for i, object := range infolist {
		tr := tb.Body().Tr()
		tr.Body().Th().Attr("scope", "row").Span().Class("badge rounded-pill text-bg-light").Text(strconv.Itoa(i + 1))
		tr.Body().Td().A().Href(serverInfo.Bucket + "/d/" + object.Key).Target("_black").Text(object.Key)
		td := tr.Body().Td()
		td.Body().Span().Class("badge text-bg-primary").Text(getObjectSizeSuitableUnit(object.Size))
	}
	divframe.Body().Hr()
	footerdiv := divframe.Body().Div().Class("container-sm text-center").Div().Class("row justify-content-sm-center").Div().Class("col-md-6")
	ul := footerdiv.Body().Ul()
	if serverInfo.Options.BeianMiit != "" {
		ul.Class("list-group list-group-horizontal")
	} else {
		ul.Class("list-group")
	}
	leftli := ul.Body().A().Class("list-group-item list-group-item-action list-group-item-light").Href("https://github.com/xinjiajuan/RollShow-go").Target("_black").Text("Powered by")
	leftli.Body().Span().Class("badge rounded-pill text-bg-success").Text("RollShow " + Config.Version)
	if serverInfo.Options.BeianMiit != "" {
		ul.Body().A().Class("list-group-item list-group-item-action list-group-item-light").Href("https://beian.miit.gov.cn/").Target("_black").Span().Class("badge text-bg-danger").Text(serverInfo.Options.BeianMiit)
	}
	divframe.Body().Br()
	return bootstrap.String()
}
