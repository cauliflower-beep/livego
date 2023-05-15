## 基础知识

开始livego项目之前，有些必备的基础需要搞清楚。

> 启动rtmps服务器时候的证书和密钥是什么？

启动rtmps服务器时，先从配置文件中获取了证书和密钥的路径，然后使用 tls.LoadX509KeyPair 函数加载证书和密钥。

证书和密钥是用于加密和解密数据传输的，通常以文件的形式存储在磁盘上。

证书和密钥的`生成`通常需要使用`证书颁发机构（CA）`或`自签名证书`。如果你使用的是自签名证书，那么你需要手动创建证书和密钥。你可以使用 OpenSSL 工具来生成自签名证书和密钥。具体的生成方法可以参考 OpenSSL 的官方文档。如果你使用的是 CA 颁发的证书，那么你需要向 CA 申请证书和密钥。申请的具体方法可以参考 CA 的官方文档。

> hlsServer与rtmpServer是什么关系？

hlsServer 和 rtmpServer 是两个不同的服务器。

hlsServer 监听指定的地址和端口，接收客户端的 HLS 请求，并将 RTMP 流转换为 HLS 流进行传输。rtmpServer 监听指定的地址和端口，接收客户端的 RTMP 请求，并将音视频流传输给客户端。

在 livego 代码中，如果 hlsServer 不为 nil，则会将 hlsServer 作为参数传递给 rtmpServer，这样 rtmpServer 就可以将 RTMP 流转换为 HLS 流进行传输。如果 hlsServer 为 nil，则说明 HLS 服务器被禁用，此时 rtmpServer 只能将音视频流传输给客户端，不能进行转换。

> RTMP流和HLS流有什么区别？二者不都是二进制吗？

RTMP 和 HLS 都是用于`音视频流传输的协议`，但它们之间存在一些区别。

`RTMP 是基于 TCP 的应用层协议`，它将音视频流封装成二进制数据包进行传输。RTMP 支持实时传输，延迟较低，但对网络带宽和稳定性要求较高。

`HLS 是基于 HTTP 的协议`，它将音视频流分割成若干个小的 TS 文件进行传输。HLS 支持点播和直播，延迟较高，但对网络带宽和稳定性要求较低。在实际应用中，RTMP 通常用于直播场景，而 HLS 通常用于点播场景。二者都是二进制数据，但是它们的传输方式和特点不同。

## 关于livego

livego是一个简单高效的直播服务器，具备以下特点：
- 安装和使用非常简单；
- 纯 Golang 编写，性能高，跨平台；
- 支持常用的传输协议、文件格式、编码格式；

后台使用 RTMP 协议接收来自客户端的音视频流，并将其转发到 HLS 和 HTTP-FLV 协议。

#### 支持的传输协议
- RTMP
- AMF
- HLS
- HTTP-FLV

#### 支持的容器格式
- FLV
- TS

#### 支持的编码格式
- H264
- AAC
- MP3

## 安装
直接下载编译好的[二进制文件](https://livego/releases)后，在命令行中执行。

#### 从 Docker 启动
执行`docker run -p 1935:1935 -p 7001:7001 -p 7002:7002 -p 8090:8090 -d gwuhaolin/livego`启动

#### 从源码编译
1. 下载源码 `git clone https://livego.git`
2. 去 livego 目录中 执行 `go build`

## 使用
1. 启动服务：执行 `livego` 二进制文件启动 livego 服务；
2. 访问 `http://localhost:8090/control/get?room=movie` 获取一个房间的 channelkey(channelkey用于推流，movie用于播放).
3. 推流: 通过`RTMP`协议推送视频流到地址 `rtmp://localhost:1935/{appname}/{channelkey}` (appname默认是`live`), 例如： 使用 `ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/{appname}/{channelkey}` 推流([下载demo flv](https://s3plus.meituan.net/v1/mss_7e425c4d9dcb4bb4918bbfa2779e6de1/mpack/default/demo.flv));
4. 播放: 支持多种播放协议，播放地址如下:
    - `RTMP`:`rtmp://localhost:1935/{appname}/movie`
    - `FLV`:`http://127.0.0.1:7001/{appname}/movie.flv`
    - `HLS`:`http://127.0.0.1:7002/{appname}/movie.m3u8`

所有配置项: 
```bash
./livego  -h
Usage of ./livego:
      --api_addr string       HTTP管理访问监听地址 (default ":8090")
      --config_file string    配置文件路径 (默认 "livego.yaml")
      --flv_dir string        输出的 flv 文件路径 flvDir/APP/KEY_TIME.flv (默认 "tmp")
      --gop_num int           gop 数量 (default 1)
      --hls_addr string       HLS 服务监听地址 (默认 ":7002")
      --hls_keep_after_end    Maintains the HLS after the stream ends
      --httpflv_addr string   HTTP-FLV server listen address (默认 ":7001")
      --level string          日志等级 (默认 "info")
      --read_timeout int      读超时时间 (默认 10)
      --rtmp_addr string      RTMP 服务监听地址 (默认 ":1935")
      --write_timeout int     写超时时间 (默认 10)
```

### [和 flv.js 搭配使用](https://github.com/gwuhaolin/blog/issues/3)

对Golang感兴趣？请看[Golang 中文学习资料汇总](http://go.wuhaolin.cn/)

