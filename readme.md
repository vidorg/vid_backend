# vid_backend

### Environment
+ `Golang` 1.11.5

### Documents
+ See [api.md](https://github.com/vidorg/Vid_Backend/tree/master/docs/api.md)
+ Run following code to generate the swagger 
[yaml doc](https://github.com/vidorg/vid_backend/blob/master/docs/api.yaml) and 
[html doc](https://github.com/vidorg/vid_backend/blob/master/docs/api.html)
```bash
sh genapi.sh

# yaml:
# python ./docs/parse_yaml.py main.go ./docs/api.yaml

# html:
# python ./docs/to_html.py ./docs/api.yaml ./docs/api.html
```

### Run
```bash
cd vid/backend
go run main.go
```

### Dependencies
+ [gin](https://github.com/gin-gonic/gin)
+ [gorm](https://github.com/jinzhu/gorm)
+ [jwt-go](https://github.com/dgrijalva/jwt-go)
+ [yaml.v2](https://github.com/go-yaml/yaml)
+ [linkedhashmap](https://github.com/emirpasic/gods)
+ [swag](https://github.com/swaggo/swag)