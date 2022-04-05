module github.com/dueckminor/mypi-tools

go 1.13

replace github.com/docker/docker/internal/testutil => gotest.tools/v3 v3.0.0

require (
	docker.io/go-docker v1.0.0
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/containerd/containerd v1.6.2 // indirect
	github.com/creack/pty v1.1.18
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/docker/internal/testutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.4
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.7
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/googollee/go-engine.io v1.4.3-0.20200220091802-9b2ab104b298 // indirect
	github.com/googollee/go-socket.io v1.6.1
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/onsi/gomega v1.19.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.3.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/ugorji/go v1.2.7 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
	golang.org/x/net v0.0.0-20220403103023-749bd193bc2b
	golang.org/x/sys v0.0.0-20220405052023-b1e9470b6e64 // indirect
	golang.org/x/text v0.3.7
	google.golang.org/genproto v0.0.0-20220401170504-314d38edb7de // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	howett.net/plist v1.0.0
)
