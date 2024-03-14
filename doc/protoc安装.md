protoc 安装
安装成功后，可以用protoactor例子中的proto生成，比对一下
确保安装成功

cell\modules\protoactor-go-0.2.0\_examples\remote-chat\messages
sh可以改成bat执行一下，如果差异过大，那就是版本不对

安装参看
https://segmentfault.com/a/1190000039767770
也就是说 go get -u google.golang.org\protobuf

注，如果之前有安装，可能不会编译更新protoc-gen-go，
进去手工编译一下
cd {GOPATH}\pkg\mod\google.golang.org\protobuf@v1.28.0\cmd\protoc-gen-go
go build
go install

即可
可以进{GOPATH}\bin看看exe是否更新