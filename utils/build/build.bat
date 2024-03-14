set PACKAGE=github.com/dfklegend/cell2/utils
set VERSION=v1.0.0
set TIME=20230823
set TAR=.\

go build -ldflags "-X %PACKAGE%/build.Version=%VERSION% -X %PACKAGE%/build.Time=%TIME%"