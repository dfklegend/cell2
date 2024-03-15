module client

go 1.18

replace github.com/dfklegend/cell2/utils => ../../../cell2/utils

replace github.com/dfklegend/cell2/apimapper => ../../../cell2/apimapper

replace github.com/dfklegend/cell2/pomelonet => ../../../cell2/pomelonet

replace github.com/dfklegend/cell2/pomeloclient => ../../../cell2/pomeloclient

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/asynkron/goconsole v0.0.0-20160504192649-bfa12eebf716
	github.com/dfklegend/cell2/pomeloclient v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell2/pomelonet v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell2/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/dfklegend/cell2/apimapper v0.0.0-00010101000000-000000000000 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/petermattis/goid v0.0.0-20221215004737-a150e88a970d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)
