# vid_backend

### Environment
+ `Golang 1.13.5 windows/amd64`

### Documents
+ Run the following code to generate the swagger api document
+ See
[api.md](https://github.com/vidorg/Vid_Backend/tree/master/docs/api.md)
[api.yaml](https://github.com/vidorg/vid_backend/blob/master/docs/api.yaml) 
[api.html](https://github.com/vidorg/vid_backend/blob/master/docs/api.html)

```bash
sh gendoc.sh
```

### Run

```bash
# cd vid_backend
go run main.go
```

### Dependencies
+ [ahlib](https://github.com/Aoi-hosizora/ahlib)
+ [ahlib-gin-gorm](https://github.com/Aoi-hosizora/ahlib-gin-gorm)
+ [gin](https://github.com/gin-gonic/gin)
+ [gorm](https://github.com/jinzhu/gorm)
+ [jwt-go](https://github.com/dgrijalva/jwt-go)
+ [yaml.v2](https://github.com/go-yaml/yaml)
+ [swag](https://github.com/swaggo/swag)
