# vdisk_monitor

##Set etcd/client build enviorenment

fix ** *cannot find package "golang.org/x/net/context"* ** problem:
```
mkdir -p $GOPATH/src/github.com/golang
cd $GOPATH/src/github.com/golang
git clone git@github.com:golang/net.git
```
##Notes & References:

### 编码手册及规范

[effective go](https://golang.org/doc/effective_go.html)

[go classroom](https://www.kancloud.cn/digest/batu-go/153540)

[台湾人写的go教程](https://polor10101.gitbooks.io/golang_note/content/goroutine.html)

[Google go 代码规范](https://github.com/golang/go/wiki/CodeReviewComments)

[Go Q&As](https://golang.org/doc/faq#goroutines)

### Etcd

[etcdctl README](https://github.com/coreos/etcd/blob/master/etcdctl/READMEv2.md)

[etcd client](https://github.com/coreos/etcd/tree/master/client)

### Channel & Select & Interface

[goroutine channel select](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/02.7.md)

[non-blocking-channel & select](https://gobyexample.com/non-blocking-channel-operations)

[go channel](http://hustcat.github.io/channel/)

[how-to-use-interfaces-in-go](http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)

### Goroutine

[goroutine工作原理](https://www.zhihu.com/question/20862617)

### Timer & Ticker

[timer](https://gobyexample.com/timers)

[tickers](https://gobyexample.com/tickers)

### Others

[go socket](http://blog.csdn.net/ahlxt123/article/details/47320161)

[go 使用共享内存](http://studygolang.com/articles/743)

[go lldp package](https://godoc.org/github.com/mdlayher/lldp)

[go lldp code](https://github.com/mdlayher/lldp)

[laws-of-reflection](https://blog.golang.org/laws-of-reflection)

[go RPC](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/8.4.md)

[Gob - go binary](https://mikespook.com/2011/03/%E3%80%90%E7%BF%BB%E8%AF%91%E3%80%91gob-%E7%9A%84%E6%95%B0%E6%8D%AE/)

## 编译到 linux 64bit
$ GOOS=linux GOARCH=amd64 go build

## GDB GO
go build -ldflags "-w"

[Go gdb](http://blog.studygolang.com/2012/12/gdb%E8%B0%83%E8%AF%95go%E7%A8%8B%E5%BA%8F/)

