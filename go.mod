module github.com/VulpesFerrilata/auth

go 1.16

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/VulpesFerrilata/go-micro/plugins/client/grpc v0.0.0-20210419083936-3c392b08f80a
	github.com/VulpesFerrilata/go-micro/plugins/server/grpc v0.0.0-20210419083936-3c392b08f80a
	github.com/VulpesFerrilata/grpc v0.0.0-20210427071000-4aed833515e5
	github.com/VulpesFerrilata/library v0.0.0-20210426114214-2fbe619c691a
	github.com/andybalholm/brotli v1.0.1-0.20200619015827-c3da72aa01ed // indirect
	github.com/asim/go-micro/plugins/server/http/v3 v3.0.0-20210408173139-0d57213d3f5c
	github.com/asim/go-micro/v3 v3.5.0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.5
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/schema v0.0.2 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kataras/golog v0.0.18 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/kataras/neffos v0.0.16 // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/mediocregopher/radix/v3 v3.5.2 // indirect
	github.com/micro/cli/v2 v2.1.2
	github.com/microcosm-cc/bluemonday v1.0.3 // indirect
	github.com/nats-io/nats-server/v2 v2.1.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/gorm v1.21.6
)
