package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type Builder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *Builder {
	return &Builder{
		loggerFunc: fn,
	}
}

func (b *Builder) AllowReqBody() *Builder {
	b.allowReqBody = true
	return b
}
func (b *Builder) AllowRespBody() *Builder {
	b.allowRespBody = true
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	start := time.Now()
	return func(ctx *gin.Context) {
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    ctx.Request.URL.String(),
		}

		if b.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			// body 读取后，需要重新赋值给 Request.Body 是一个readCloser stream对象
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			al.ReqBody = string(body)
		}
		if b.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             al,
			}
		}

		ctx.Next()

		defer func() {
			al.Duration = time.Since(start).String()
			al.Status = ctx.Writer.Status()
			b.loggerFunc(ctx, al)
		}()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}
func (w *responseWriter) WriteString(s string) (int, error) {
	w.al.RespBody = s
	return w.ResponseWriter.WriteString(s)
}
func (w *responseWriter) WriteHeader(code int) {
	w.al.Status = code
	w.ResponseWriter.WriteHeader(code)
}

type AccessLog struct {
	// HTTP 方法
	Method string
	// 请求路径
	Url string

	Status   int
	ReqBody  string
	RespBody string
	Duration string
}
