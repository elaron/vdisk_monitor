# vdisk_monitor

Set etcd/client build enviorenment

fix ** *cannot find package "golang.org/x/net/context"* ** problem:
```
mkdir -p $GOPATH/src/github.com/golang
cd $GOPATH/src/github.com/golang
git clone git@github.com:golang/net.git
```
Notes & References:

[goroutine channel](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/02.7.md)

[go socket](http://blog.csdn.net/ahlxt123/article/details/47320161)

[etcdctl README](https://github.com/coreos/etcd/blob/master/etcdctl/READMEv2.md)

[etcd client](https://github.com/coreos/etcd/tree/master/client)
[laws-of-reflection](https://blog.golang.org/laws-of-reflection)

# 编译到 linux 64bit
$ GOOS=linux GOARCH=amd64 go build

# GDB GO
go build -ldflags "-w"

