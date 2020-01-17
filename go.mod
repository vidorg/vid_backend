module github.com/vidorg/vid_backend

require (
	github.com/Aoi-hosizora/ahlib v0.0.0-20200117141519-e2c13dbee158
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.3 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.10
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.3
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b
	golang.org/x/net v0.0.0-20191011234655-491137f69257 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/yaml.v2 v2.2.4
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.50.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191227163750-53104e6ec876
	golang.org/x/exp => github.com/golang/exp v0.0.0-20191227195350-da58074b4299
	golang.org/x/image => github.com/golang/image v0.0.0-20191214001246-9130b4cfad52
	golang.org/x/net => github.com/golang/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20200102141924-c96a22e43c9c
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200103221440-774c71fcf114
	google.golang.org/appengine => github.com/golang/appengine v1.6.5
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20191230161307-f3c370f40bfb
	google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0
)

go 1.13
