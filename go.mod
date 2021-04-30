module github.com/dueckminor/mypi-tools

go 1.13

replace github.com/docker/docker/internal/testutil => gotest.tools/v3 v3.0.0

require (
	docker.io/go-docker v1.0.0
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/creack/pty v1.1.11
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/docker/internal/testutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fatih/color v1.9.0
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-contrib/static v0.0.0-20191128031702-f81c604d8ac2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/googollee/go-socket.io v1.4.3
	github.com/gorilla/websocket v1.4.2
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.2.0
	github.com/shirou/gopsutil v3.21.2+incompatible
	github.com/skip2/go-qrcode v0.0.0-20200526175731-7ac0b40b2038
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	golang.org/x/text v0.3.2
	gopkg.in/yaml.v2 v2.3.0
	howett.net/plist v0.0.0-20200419221736-3b63eb3a43b5
)
