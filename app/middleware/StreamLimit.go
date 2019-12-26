package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"vid/app/controller/exception"
	"vid/app/model/dto"
)

func StreamLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		rawBuf, err := c.GetRawData()
		if err == nil {
			// [GIN-debug] error on parse multipart form array: multipart: NextPart: EOF
			buf := bytes.NewBuffer(rawBuf)
			c.Request.Body = ioutil.NopCloser(buf)
			c.Next()
		} else {
			conn, bufRw, err := c.Writer.Hijack()
			if err == nil {
				_ = bufRw.Flush()
				_ = conn.Close()
			}
			c.JSON(http.StatusRequestEntityTooLarge, dto.Result{}.Error(http.StatusRequestEntityTooLarge).SetMessage(exception.RequestSizeLargeError.Error()))
			c.Abort()
		}
	}
}
