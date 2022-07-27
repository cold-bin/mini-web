// @author cold bin
// @date 2022/7/26

package mini_web

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type H map[string]interface{}

// Context 将一次请求和应答的所有东西塞入上下文处理，上下文统一提供方便的接口，便于编写优雅业务代码
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	engine *Engine

	//context map
	Keys map[string]interface{}

	// read & write mutex for the Keys
	Mu sync.RWMutex
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
		Keys:   map[string]interface{}{},
	}
}

func (c *Context) ShouldBind(obj interface{}) error {

	switch c.GetHeader("Content-Type") {
	case "application/json":
		return c.ShouldBindJson(obj)
	case "application/xml":
		return c.ShouldBindXml(obj)
	default:
		return errors.New("not support the request Content-Type")
	}
}

func (c *Context) ShouldBindJson(obj interface{}) error {
	//GET have not the body
	if c.Req.Method == http.MethodGet {
		return errors.New(http.MethodGet + "have no the body")
	}
	size := c.GetHeader("Content-length")
	bytesNum, err := strconv.Atoi(size)
	if err != nil {
		return err
	}

	bytes := make([]byte, bytesNum)
	reader := bufio.NewReader(c.Req.Body)
	if _, err := reader.Read(bytes); err != nil {
		return err
	}

	return json.Unmarshal(bytes, obj)
}

func (c *Context) ShouldBindXml(obj interface{}) error {
	//GET have not the body
	if c.Req.Method == http.MethodGet {
		return errors.New(http.MethodGet + "have no the body")
	}

	size := c.GetHeader("Content-length")
	bytesNum, err := strconv.Atoi(size)
	if err != nil {
		return err
	}
	bytes := make([]byte, bytesNum)

	return xml.Unmarshal(bytes, obj)
}

func (c *Context) Set(key string, val string) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	// lock the key
	c.Mu.Lock()
	c.Keys[key] = val
	c.Mu.Unlock()
}

func (c *Context) Get(key string) (interface{}, bool) {
	c.Mu.RLock()
	val, ok := c.Keys[key]
	c.Mu.RUnlock()
	return val, ok
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) DefaultPostForm(key, defaultStr string) string {
	if val := c.PostForm(key); val != "" {
		return val
	}
	return defaultStr
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) DefaultQuery(key, defaultStr string) string {
	if val := c.Query(key); val != "" {
		return val
	}

	return defaultStr
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.String(), nil
}

func (c *Context) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	if path == "" {
		path = "/"
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
