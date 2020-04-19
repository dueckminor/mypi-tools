module github.com/dueckminor/mypi-tools

go 1.13

replace github.com/docker/docker/internal/testutil => gotest.tools/v3 v3.0.0

require (
	docker.io/go-docker v1.0.0
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/creack/pty v1.1.9
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/docker/internal/testutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fatih/color v1.9.0
	github.com/fatih/structs v1.1.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-contrib/static v0.0.0-20191128031702-f81c604d8ac2
	github.com/gin-gonic/gin v1.6.2
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/googollee/go-socket.io v1.4.3
	github.com/gorilla/websocket v1.4.2
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.2.0
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	gopkg.in/yaml.v2 v2.2.8
	howett.net/plist v0.0.0-20200225050739-77e249a2e2ba
)
