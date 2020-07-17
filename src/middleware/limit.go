package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
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
			result.Error(exception.RequestLargeError).JSON(c)
			c.Abort()
			return
		}

		// [GIN-debug] error on parse multipart form array: multipart: NextPart: EOF
		buf := bytes.NewBuffer(rawBuf)
		c.Request.Body = ioutil.NopCloser(buf)
		c.Next()
	}
}
