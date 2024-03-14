set PACKAGE=github.com/dfklegend/cell2/utils
set VERSION=v1.0.0
set TIME=20230823
set TAR=.\main

go build %TAR% -ldflags="-X 'main.Version=1.0.0'"