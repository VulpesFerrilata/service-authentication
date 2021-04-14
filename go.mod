module github.com/VulpesFerrilata/auth

go 1.14

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/VulpesFerrilata/grpc v0.0.0-20210408081122-0a9eabd471f9
	github.com/VulpesFerrilata/library v0.0.0-20210414093404-58886934a95b
	github.com/andybalholm/brotli v1.0.1-0.20200619015827-c3da72aa01ed // indirect
	github.com/asim/go-micro/v3 v3.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.5.0
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/schema v0.0.2 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/golog v0.0.18 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/kataras/neffos v0.0.16 // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mediocregopher/radix/v3 v3.5.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/microcosm-cc/bluemonday v1.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/driver/mysql v1.0.1
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.21.6
)
