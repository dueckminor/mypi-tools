module github.com/dueckminor/mypi-tools

go 1.21

toolchain go1.22.4

replace github.com/docker/docker/internal/testutil => gotest.tools/v3 v3.5.1

replace github.com/docker/distribution v2.8.1+incompatible => github.com/docker/distribution v2.8.3+incompatible

require (
	docker.io/go-docker v1.0.0
	github.com/bhendo/go-powershell v0.0.0-20190719160123-219e7fb4e41e
	github.com/creack/pty v1.1.21
	github.com/docker/docker v27.1.1+incompatible
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/fatih/color v1.17.0
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-contrib/sessions v1.0.1
	github.com/gin-contrib/static v1.1.2
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/googollee/go-socket.io v1.7.0
	github.com/gorilla/websocket v1.5.3
	github.com/influxdata/influxdb-client-go/v2 v2.13.0
	github.com/miekg/dns v1.1.61
	github.com/onsi/gomega v1.33.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.6
	github.com/pquerna/otp v1.4.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/simonvetter/modbus v1.6.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	golang.org/x/crypto v0.31.0
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8
	golang.org/x/net v0.33.0
	golang.org/x/term v0.29.0
	golang.org/x/text v0.21.0
	gopkg.in/yaml.v3 v3.0.1
	howett.net/plist v1.0.1
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker/internal/testutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goburrow/serial v0.1.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.2.2 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/juju/errors v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/oapi-codegen/runtime v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/quasoft/memstore v0.0.0-20191010062613-2bce066d2b0b // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel v1.25.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.25.0 // indirect
	go.opentelemetry.io/otel/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/sdk v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)
