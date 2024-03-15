module mmo

go 1.21

toolchain go1.22.1

replace (
	github.com/dfklegend/cell2 => ./../../../
	github.com/dfklegend/cell2/apimapper => ./../../../apimapper
	github.com/dfklegend/cell2/pomelonet => ./../../../pomelonet
	github.com/dfklegend/cell2/utils => ./../../../utils
)

require (
	github.com/asynkron/protoactor-go v0.0.0-20240308120642-ef91a6abee75
	github.com/dfklegend/astar v0.0.0-20231101082850-5c38c64d9c8c
	github.com/dfklegend/cell2 v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell2/apimapper v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell2/utils v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v7 v7.4.1
	github.com/magiconair/properties v1.8.7
	github.com/spf13/viper v1.14.0
	github.com/stretchr/testify v1.8.4
	github.com/yuin/gopher-lua v1.1.0
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	google.golang.org/protobuf v1.31.0
	layeh.com/gopher-luar v1.0.10
)

require (
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/Workiva/go-datastructures v1.1.1 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/asynkron/gofun v0.0.0-20220329210725-34fed760f4c2 // indirect
	github.com/aws/aws-sdk-go v1.33.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cbroglie/mustache v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cheggaaa/pb/v3 v3.0.5 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dfklegend/cell2/pomelonet v0.0.0-00010101000000-000000000000 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gocarina/gocsv v0.0.0-20221212090504-363f27c9ceb6 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lib/pq v1.7.0 // indirect
	github.com/lithammer/shortuuid/v4 v4.0.0 // indirect
	github.com/lmittmann/tint v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.6.3 // indirect
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/petermattis/goid v0.0.0-20221215004737-a150e88a970d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/vadv/gopher-lua-libs v0.4.1 // indirect
	github.com/yuin/gluamapper v0.0.0-20150323120927-d836955830e7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.10 // indirect
	go.etcd.io/etcd/client/v3 v3.5.10 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
