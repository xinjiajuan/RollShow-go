# RollShow
## 支持S3对象存储的文件下载网站服务器
用goland编写，主要功能是可以连接S3对象存储服务器，并且监听端口将配置的存储服务器中桶的文件通过网页展示出来，并且S3服务器连接只在与本服务器，通过路由断绝，用户端看不到S3服务器后端链接。
一个服务器，一个桶，一个协程，互不干扰

# 为什么要做这样一个东西？
minio是用的第一个对象存储服务器，现在已经用来做我很多资料的存储和备份。以至于博客的资源文件也放到对象存储中。但是问题也随之而来：博客资源主要使用直链进行资源的调用，而minio提供了api也是支持直链访问，但是桶需要改为开放类型，而允许所有人对其读取和修改，但是改为私人桶，访问直链又需要鉴权，可博客没有对于s3鉴权的功能。对于访客来说修改存储的资源是不允许的，所以这个项目由此而生，来展示支持s3对象存储的文件，可以生成直链，且实时与对象存储后端保持同步，而且保存在存储服务器中的文件也不会被修改，很安全。



# 启动
```shell
$ .\rollshow.exe -c .\Config.yaml #windows&linux
#记录log并后台运行
$ ./rollshow.exe -c config.yaml >> rollshow.log & #linux
```

# 从源码编译


```shell
# 拉取源码
$ git clone https://github.com/xinjiajuan/RollShow-go.git
# 进入源码文件夹
$ cd RollShow-go
# 拉取软件需要的包
$ go mod tidy
# 编译
$ go build
# 给二进制执行权限
$ chmod +x rollshow #linux
# 运行
$ ./rollshow -c config.yaml
```
# 配置文件
程序启动需要指定配置文件,可使用绝对路径和相对路径。

每一个`- name`都是一个实例，互不干扰。

```yaml
server:
  - name: minio1 #名称，用于方便用户标识实例，无实际意义，必填
    listenPort: 8080 #监听端口，必填
    enable: true #是否启用此实例，必填
    host: 192.168.2.220:9000 #S3API的Url，必填
    accessKeyID: 'API user' #顾名思义，必填
    secretAccessKey: qwertyuio #顾名思义，必填
    bucket: blog-res #桶，必填
    options: #其他参数,不填
      useSSL: true # 启用TLS连接s3服务器，必填
      region: chinaxxxxxx #区域，选填
      bucketLookupType: 0 #桶查找类型 DNS:2,Path:1,Auto:0，必填
    web: #web相关设置
      useTLS: #web网页tls设置
        enable: false
        certFile: '/home/user/cert.pem'
        certKey: '/home/user/key.pem'
      access-control-allow-origin: 'www.domain.com' #资源跨域策略,只对下载链接有效,主页无跨域设置
      favicon: "blog-res/d/blog/ico_s/logo.png" #网页图标url,通过 301 跳转获取,暂不支持本地图片,请使用在线资源
      beianMiit: "" #工信部的备案号,显示在前端，为空不显示
  - name: minio2
    listenPort: 8081
    enable: true
    host: 192.168.2.220:9000
    accessKeyID: user
    secretAccessKey: xxxxxx
    bucket: bucket
    options: #其他参数,不填
      useSSL: true # 启用TLS连接s3服务器，必填
      region: chinaxxxxxx #区域，选填
      bucketLookupType: 0 #桶查找类型 DNS:2,Path:1,Auto:0，必填
    web: #web相关设置
      useTLS: #web网页tls设置
        enable: false
        certFile: '/home/user/cert.pem'
        certKey: '/home/user/key.pem'
      access-control-allow-origin: 'www.domain.com' #资源跨域策略,只对下载链接有效,主页无跨域设置
      favicon: "blog-res/d/blog/ico_s/logo.png" #网页图标url,通过 301 跳转获取,暂不支持本地图片,请使用在线资源
      beianMiit: "" #工信部的备案号,显示在前端，为空不显示
```

# 第三方库

- github.com/klarkxy/gohtml
- github.com/minio/minio-go/
- github.com/urfave/cli/
- gopkg.in/yaml.v3

# 性能展示
下面性能展示使用的服务器配置为：
- 服务器型号：Inspur NF5280M3
- CPU: Intel E5-2650V2
- RAM: 8G ECC
- 硬盘: HGST MSIP-REM-HG2-HUC101890CSS20
- 阵列卡: E300750 单盘 Raid 0

图中`192.168.2.220`主机为上诉测试服务器

## 首页性能
> 局域网1G带宽下1000并发

![image](https://user-images.githubusercontent.com/36360150/181248990-7bff889a-1ec7-4f85-8958-cb607ad6f081.png)

## 下载性能
> 局域网1G带宽下100并发下载32M无损歌曲

![image](https://user-images.githubusercontent.com/36360150/181250742-d76f904b-7741-4ad4-9bbc-c9b2551be90e.png)

# 日志
## 2022-10-22 v1.2.0

- 优化部分配置字段
- 添加web支持tls，web页面支持https协议

## 2022-8-16 v1.1.8

- 支持配置`/d/`路径下的跨域配置
- 支持显示`favicon`

## 2022-7-30 v1.1.5 beta

- 支持多线程和断点续传下载
- 优化服务器 log 打印显示

## 2022-7-27 v1.1.1 beta

- 优化文件大小单位的计算方法
- 优化前端界面
- 修复了连接s3后端时出现错误不显示直接程序报错的问题

## 2022-7-27 v1.1.0 beta

- 项目改名为`RollShow`,中文名为`展卷`,意为将卷轴展开，代表程序的最本质功能是为s3对象存储桶对象展示与下载
- 程序正式可用，前端与后端下载功能Debug正常，可以正式使用；但服务器性能有待评价
- 前端ui调整完毕

## 2022-7-25 v1.0.5 Debug

- 使用bootstrap渲染前端
- 使用gohtml生成html页面，库地址: https://github.com/klarkxy/gohtml
- 前端支持显示文件大小、文件在S3中的路径、文件序号

## 2022-7-19 v1.0.2 Debug

- 使用`Ctrl+c`关闭程序时将平滑关闭服务实例再退出
## 2022-7-18 v1.0.1 Debug

- S3对象信息能成功读取，并且能将信息成功打印至监听网页
- YAML配置文件已经完善
- 程序启动使用控制台，带-c参数
