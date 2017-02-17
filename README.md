# vdisk_monitor

Set etcd/client build enviorenment

fix ** *cannot find package "golang.org/x/net/context"* ** problem:
```
mkdir -p $GOPATH/src/github.com/golang
cd $GOPATH/src/github.com/golang
git clone git@github.com:golang/net.git
```
Notes & References:

[Google go 代码规范](https://github.com/golang/go/wiki/CodeReviewComments)

[goroutine channel select](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/02.7.md)

[goroutine工作原理](https://www.zhihu.com/question/20862617)

[non-blocking-channel & select](https://gobyexample.com/non-blocking-channel-operations)

[timer](https://gobyexample.com/timers)

[go socket](http://blog.csdn.net/ahlxt123/article/details/47320161)

[etcdctl README](https://github.com/coreos/etcd/blob/master/etcdctl/READMEv2.md)

[etcd client](https://github.com/coreos/etcd/tree/master/client)

[laws-of-reflection](https://blog.golang.org/laws-of-reflection)

[go classroom](https://www.kancloud.cn/digest/batu-go/153540)

[go RPC](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/8.4.md)

[Gob](https://mikespook.com/2011/03/%E3%80%90%E7%BF%BB%E8%AF%91%E3%80%91gob-%E7%9A%84%E6%95%B0%E6%8D%AE/)

[Go gdb](http://blog.studygolang.com/2012/12/gdb%E8%B0%83%E8%AF%95go%E7%A8%8B%E5%BA%8F/)

[台湾人写的go教程](https://polor10101.gitbooks.io/golang_note/content/goroutine.html)

[Go Q&As](https://golang.org/doc/faq#goroutines)

# 编译到 linux 64bit
$ GOOS=linux GOARCH=amd64 go build

# GDB GO
go build -ldflags "-w"

