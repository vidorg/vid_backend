package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/model/common"
	"io/ioutil"
	"net/http"
)

func LimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		rawBuf, err := c.GetRawData()
		if err != nil {
			// https://github.com/gin-gonic/gin/issues/1136
			// conn, bufRw, err := c.Writer.Hijack()
			common.Result{}.Error(http.StatusRequestEntityTooLarge).SetMessage(exception.RequestSizeError.Error()).JSON(c)
			c.Abort()
			return
		}

		// [GIN-debug] error on parse multipart form array: multipart: NextPart: EOF
		buf := bytes.NewBuffer(rawBuf)
		c.Request.Body = ioutil.NopCloser(buf)
		c.Next()
	}
}
