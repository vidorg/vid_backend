# vid_backend

### Environment
+ `Golang 1.13.5 windows/amd64`

### Documents
+ Run the following code to generate the swagger api document
+ See
[api.yaml](https://github.com/vidorg/vid_backend/blob/master/docs/api.yaml) 
[api.html](https://github.com/vidorg/vid_backend/blob/master/docs/api.html)

```bash
# To generate api document
sh gendoc.sh
```

### Run

```bash
# To run directly
go run main.go

# To build
go build -i -o ./build/vid_backend.out main.go
./build/vid_backend.out
```

### Dependencies
+ [ahlib](https://github.com/Aoi-hosizora/ahlib)
+ [ahlib-gin-gorm](https://github.com/Aoi-hosizora/ahlib-gin-gorm)
+ [gin](https://github.com/gin-gonic/gin)
+ [gorm](https://github.com/jinzhu/gorm)
+ [jwt-go](https://github.com/dgrijalva/jwt-go)
+ [yaml.v2](https://github.com/go-yaml/yaml)
+ [swag](https://github.com/swaggo/swag)
+ [logrus](https://github.com/sirupsen/logrus)
