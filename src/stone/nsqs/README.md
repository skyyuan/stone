# NSQ Server

这是基于[nsq](http://nsq.io)的一个封装，能方便的基于nsq构筑一个服务系统。

## 使用

```go
import "gitlab.chainresearch.org/wallet/stone/nsqs"

config := &nsqs.SimpleConfig{}

// 初始化
nsqs.InitConfig(conf)

// 注册监听
nsqs.Register("test_topic", "test_channel", testFunc, 10)

// 运行server
nsqs.Run()
```

其中`testFunc`为一个符合`func (m *nsq.Message) error`的方法。

方法返回error的时候或者panic的时候会自动重试。

### 发送消息

```go
nsqs.GlobalEmmiter.Emit("topic", payload)
```

### 与其他服务框架共用的时候

例如，与`echo`共用

```go
e := echo.New()

// 启动nsqs
nsqs.Start()

// 等待echo退出
quit := make(chan os.Signal)
signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 停止nsqs，并等待执行完成
nsqs.Stop()
// echo退出
e.Shutdown(ctx)
```