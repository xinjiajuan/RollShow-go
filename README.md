# S3对象存储网络文件浏览服务器
用goland编写，主要功能是可以连接S3对象存储服务器，并且监听端口将配置的存储服务器中桶的文件通过网页展示出来，并且S3服务器连接只在与本服务器，通过路由断绝，用户端看不到S3服务器后端链接。

效果就像用nginx搭建的本地文件服务器一样。
一个桶一个服务器runtime一个端口，互不干扰
# 启动
```shell
.\go_build_S3ObjectStorageFileBrowser.exe -c .\Config.yaml
格式:二进制程序 -c 配置文件路径
```
# 配置文件
可使用绝对路径和相对路径。
每一个`- name`都是一个实例，互不干扰。

```yaml
server:
  - name: minio1 #名称，用于方便用户标识实例，无实际意义
    listenPort: 8080 #监听端口
    enable: true #是否启用此实例
    host: 192.168.2.220:9000 #S3API的Url
    accessKeyID: 'API user' #顾名思义
    secretAccessKey: qwertyuio #顾名思义
    bucket: blog-res #桶
    options: #复合选项
      useSSL: true # 启用TLS连接服务器吗
      region: chinaxxxxxx #区域
      bucketLookupType: 0 #桶查找类型 DNS,Path:1,Auto:0
  - name: minio2
    listenPort: 8080
    enable: true
    host: 192.168.2.220:9000
    accessKeyID: user
    secretAccessKey: xxxxxx
    bucket: bucket
    options:
      useSSL: true
      region: china-xxxxxx
      bucketLookupType: 0 #DNS,Path:1,Auto:0
```

# 日志
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