# stone

钱包后端API。

## 开发

### 环境

* go 版本：1.8.3
* mysql版本：5.6
* 依赖管理工具：[dep](https://github.com/golang/dep)
* 单元测试工具：[testify](https://github.com/stretchr/testify)
* 消息队列：[nsq](http://nsq.io/)

### 编辑器

推荐使用[VSCode](https://code.visualstudio.com/)

### i18n

使用工具：[go-i18n](https://github.com/nicksnyder/go-i18n)

目录：`locale`

编辑各语言文件，在文件中增加对应的`key`以及翻译文本。然后执行`fileb0x b0x.yaml`将语言文件转换为go文件。

### 日志

日志分为2个部分：请求日志和调试日志。

#### 请求日志

请求日志记录每一条API请求的信息，主要用于数据分析。如果有系统级的panic，也会记录下来。

日志记录的信息如下：

* hostname：机器名
* statusCode：状态码
* latency：处理耗时
* clientIP：客户端IP
* method：请求方法
* path：请求路径
* dataLength：返回值大小
* userAgent：请求的 user agent
* appKey：app key
* headers：请求头hash
* requestID：请求ID
* errcode：错误码
* errmsg：错误消息
* panicStack：异常堆栈

该日志在框架中统一处理。使用工具： [logrus](https://github.com/Sirupsen/logrus)。

* 日志格式: `json`
* 输出: 文件，根据需要可灵活拓展为`kafka`、`logstash`等后端
    * 开发环境下为`stdout`

#### 调试日志

调试日志为开发／线上调试用的日志，不需要接入数据分析，只使用`debug`模式输出。通过启动参数`--debug`开启。

* 日志格式：文本
* 输出：`stdout`

### 消息队列

### 启动项目
go run ./cmd/{N}/main.go

