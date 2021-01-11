module github.com/VulpesFerrilata/auth

go 1.14

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/VulpesFerrilata/go-micro-custom v0.0.0-20210108071325-c08180f3a0a1
	github.com/VulpesFerrilata/grpc v0.0.0-20210108084505-964ee54318c0
	github.com/VulpesFerrilata/library v0.0.0-20210108114535-906b7f5cafa7
	github.com/andybalholm/brotli v1.0.1-0.20200619015827-c3da72aa01ed // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/universal-translator v0.17.0
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f // indirect
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/schema v0.0.2 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kataras/golog v0.0.18 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/kataras/neffos v0.0.16 // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/mediocregopher/radix/v3 v3.5.2 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/microcosm-cc/bluemonday v1.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/gorm v1.20.1
)
