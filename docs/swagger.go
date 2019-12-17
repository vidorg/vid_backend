package docs

import (
	"github.com/swaggo/swag"
	"io/ioutil"
	"os"
)

type swagger struct{}

func (s *swagger) ReadDoc() string {
	f, err := os.Open("./docs/api.yaml")
	if err != nil {
		return ""
	}
	buf, err := ioutil.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return ""
	}
	return string(buf)
}

func init() {
	swag.Register(swag.Name, &swagger{})
}
