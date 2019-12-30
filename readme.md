# vid_backend

### Environment
+ `Golang` 1.11.5

### Documents
+ Run following code to generate the swagger api document

+ See
[api.md](https://github.com/vidorg/Vid_Backend/tree/master/docs/api.md) and
[api.yaml](https://github.com/vidorg/vid_backend/blob/master/docs/api.yaml) and 
[api.html](https://github.com/vidorg/vid_backend/blob/master/docs/api.html)

```bash
sh gendoc.sh

# yaml:
# python ./api_yaml.py main.go ./docs/api.yaml

# html:
# python ./api_html.py ./docs/api.yaml ./docs/api.html

# markdown
# npm i -g swagger-markdown
# swagger-markdown -i ./docs/api.yaml -o ./docs/api.md
```

### Run

```bash
# cd vid_backend
go run main.go
```

### Dependencies
+ [gin](https://github.com/gin-gonic/gin)
+ [gorm](https://github.com/jinzhu/gorm)
+ [jwt-go](https://github.com/dgrijalva/jwt-go)
+ [yaml.v2](https://github.com/go-yaml/yaml)
+ [linkedhashmap](https://github.com/emirpasic/gods)
+ [swag](https://github.com/swaggo/swag)
