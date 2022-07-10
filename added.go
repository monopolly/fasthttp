package fasthttp

//functions to original fasthttp

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/monopolly/cast"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	ContentTypeJS            = "application/javascript"
	ContentTypeJSON          = "application/json"
	ContentTypeURL           = "application/x-www-form-urlencoded"
	ContentTypeXML           = "application/xml"
	ContentTypeZIP           = "application/zip"
	ContentTypePDF           = "application/pdf"
	ContentTypeMP3           = "audio/mpeg"
	ContentTypeVorbis        = "audio/vorbis"
	ContentTypeMultipartForm = "multipart/form-data"
	ContentTypeCSS           = "text/css"
	ContentTypeHTML          = "text/html"
	ContentTypePlain         = "text/plain"
	ContentTypePNG           = "image/png"
	ContentTypeJPG           = "image/jpeg"
	ContentTypeGIF           = "image/gif"
	ContentTypeSSE           = "text/event-stream"
	ContentTypeIcon          = "image/ico"
)

type Context = RequestCtx
type Framework = Server

func New(port int, handler RequestHandler) (a *Server) {
	a = &Server{
		Handler: handler,
		Name:    "nginx 1.1.10",
	}
	go a.ListenAndServe(fmt.Sprintf(":%d", port))
	//log.Println("starting...", port)
	return
}

func NewTLS(port int, handler RequestHandler, certfile, keyfile string) (a *Server) {
	a = &Server{
		Handler: handler,
		Name:    "nginx 1.1.10",
	}
	go a.ListenAndServeTLS(fmt.Sprintf(":%d", port), certfile, keyfile)
	log.Println("starting tls...", port)
	return
}

func (c *Context) ContentTypeJS()            { c.SetContentType(ContentTypeJS) }
func (c *Context) ContentTypeJSON()          { c.SetContentType(ContentTypeJSON) }
func (c *Context) ContentTypeXML()           { c.SetContentType(ContentTypeXML) }
func (c *Context) ContentTypeZIP()           { c.SetContentType(ContentTypeZIP) }
func (c *Context) ContentTypePDF()           { c.SetContentType(ContentTypePDF) }
func (c *Context) ContentTypeMP3()           { c.SetContentType(ContentTypeMP3) }
func (c *Context) ContentTypeVorbis()        { c.SetContentType(ContentTypeVorbis) }
func (c *Context) ContentTypeMultipartForm() { c.SetContentType(ContentTypeMultipartForm) }
func (c *Context) ContentTypeCSS()           { c.SetContentType(ContentTypeCSS) }
func (c *Context) ContentTypeHTML()          { c.SetContentType(ContentTypeHTML) }
func (c *Context) ContentTypePlain()         { c.SetContentType(ContentTypePlain) }
func (c *Context) ContentTypePNG()           { c.SetContentType(ContentTypePNG) }
func (c *Context) ContentTypeJPG()           { c.SetContentType(ContentTypeJPG) }
func (c *Context) ContentTypeGIF()           { c.SetContentType(ContentTypeGIF) }
func (c *Context) ContentTypeSSE()           { c.SetContentType(ContentTypeSSE) }
func (c *Context) ContentTypeIcon()          { c.SetContentType(ContentTypeIcon) }

func (c *Context) CSS(css []byte) {
	c.SetContentType(ContentTypeCSS)
	c.Write(css)
}

func (c *Context) CacheForever() {
	c.Response.Header.Add("Cache-Control", "max-age=31536000")
}

func (c *Context) CacheMax(sec int) {
	c.Response.Header.Add("Cache-Control", fmt.Sprintf("max-age=%d", sec))
}

func (c *Context) CacheTag(tag string) {
	c.Response.Header.Add("ETag", fmt.Sprintf(`"%s"`, tag))
}

func (c *Context) CacheNoStore() {
	c.Response.Header.Add("Cache-Control", "no-store")
}

func (c *Context) Route() (v string) {
	return fmt.Sprintf("%s %s", string(c.Method()), string(c.Path()))
}

func (c *Context) Header(key string) (v string) {
	return string(c.Request.Header.Peek(key))
}

func (c *Context) Cookie(key string) (v string) {
	return string(c.Request.Header.Cookie(key))
}

func (a *Context) NewCookie(key string, value string) (cookie *Cookie) {
	cookie = AcquireCookie()
	//cookie.SetPath("/")
	cookie.SetKey(key)
	cookie.SetValue(value)
	cookie.SetHTTPOnly(true)
	return
}

func (a *Context) AddCookie(cookie *Cookie) {
	a.Response.Header.SetCookie(cookie)
}

//если ставишь с доменом то и удалять нужно с доменом! иначе не удалится
func (a *Context) DeleteAllCookies(domain ...string) {
	a.Request.Header.VisitAllCookie(func(k, v []byte) {
		cc := a.NewCookie(string(k), "")
		if len(domain) > 0 {
			cc.SetDomain(domain[0])
		}
		cc.SetExpire(time.Now().AddDate(-20, 0, 0))
		a.AddCookie(cc)
	})
}

func (a *Context) DeleteCookie(key string, domain ...string) {
	n := AcquireCookie()
	n.SetKey(key)
	n.SetExpire(time.Now().AddDate(-10, 0, 0))
	if domain != nil {
		n.SetDomain(domain[0])
	}
	a.AddCookie(n)
}

//пример c.SetCookie("", "admin", false, token, time.Now().AddDate(0, 0, 10))
func (a *Context) SetCookie(domain, key string, secure bool, value interface{}, expired ...time.Time) {
	n := AcquireCookie()
	n.SetKey(key)
	n.SetPath("/")
	n.SetValue(fmt.Sprint(value))
	n.SetHTTPOnly(true)
	n.SetSecure(secure)
	n.SetDomain(domain)
	if expired != nil {
		n.SetExpire(expired[0])
	}
	a.AddCookie(n)
}

func (a *Context) SetCookiePublic(domain, key string, value interface{}, expired ...time.Time) {
	n := AcquireCookie()
	n.SetKey(key)
	n.SetValue(fmt.Sprint(value))
	n.SetHTTPOnly(false)
	n.SetSecure(false)
	n.SetDomain(domain)
	if len(expired) > 0 {
		n.SetExpire(expired[0])
	}
	a.AddCookie(n)
}

func (c *Context) SetUser(id int) {
	c.SetUserValue("user", id)
}

func (c *Context) User() (v int) {
	v, _ = c.UserValue("user").(int)
	return
}

func (c *Context) SetPin(pin bool) {
	c.SetUserValue("pin", pin)
}

func (c *Context) Pin() (v bool) {
	v, _ = c.UserValue("pin").(bool)
	return
}

func (c *Context) SetAdmin(v bool) {
	c.SetUserValue("admin", v)
}

func (c *Context) Admin() (v bool) {
	v, _ = c.UserValue("admin").(bool)
	return
}

//json
type js struct {
	c *Context
	k map[string]interface{}
}

func (a *js) Add(k string, v interface{}) *js {
	a.k[k] = v
	return a
}
func (a *js) Send() (c *Context) {
	b, _ := ffjson.Marshal(a.k)
	a.c.Write(b)
	return a.c
}

func (c *Context) J(k string, v interface{}) (j *js) {
	j = new(js)
	j.c = c
	j.k = map[string]interface{}{
		k: v,
	}
	return
}

func (c *Context) Json(v interface{}) {
	b, _ := ffjson.Marshal(v)
	c.Write(b)
	return
}

func (c *Context) SetProject(project []byte) {
	c.SetUserValue("project", project)
}

func (c *Context) Project() (v []byte) {
	v, _ = c.UserValue("project").([]byte)
	return
}

//устанавливает IP
func (c *Context) SetIP(k string) {
	c.SetUserValue("_ip", k)
}

//читает IP если он с nginx например
func (c *Context) IP() (v string) {
	v, _ = c.UserValue("_ip").(string)
	if v == "" {
		return c.RemoteIP().String()
	}
	return
}

func (c *Context) SetEmailConfirm(k bool) {
	c.SetUserValue("confirm", k)
}

func (c *Context) EmailConfirmed() (v bool) {
	v, _ = c.UserValue("confirm").(bool)
	return
}

func (a *Context) HTML(v []byte) {
	a.SetContentType(ContentTypeHTML)
	a.Write(v)
}

/* URL params */

func (a *Context) URLParam(key string) string {
	return string(a.Request.URI().QueryArgs().Peek(key))
}
func (a *Context) URLParamBytes(key string) []byte {
	return a.Request.URI().QueryArgs().Peek(key)
}

func (a *Context) URLParamInt(key string) int {
	return cast.Int(string(a.Request.URI().QueryArgs().Peek(key)))
}
func (a *Context) URLParamFloat(key string) float64 {
	return cast.Float(string(a.Request.URI().QueryArgs().Peek(key)))
}

func (a *Context) URLParamBool(key string) bool {
	return cast.Bool(string(a.Request.URI().QueryArgs().Peek(key)))
}

func (a *Context) URLParamSlice(key string) (r []string) {
	for _, s := range strings.Split(a.URLParam(key), ",") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			r = append(r, s)
		}
	}
	return
}

func (a *Context) URLParamSliceInt(key string) (r []int) {
	for _, s := range strings.Split(a.URLParam(key), ",") {
		s = strings.TrimSpace(s)
		e := cast.Int(s)
		if e > 0 {
			r = append(r, e)
		}
	}
	return
}

func (a *Context) Param(id string) string {
	return cast.String(a.UserValue(id))
}

func (a *Context) ParamInterface(key string) interface{} {
	return a.UserValue(key)
}

func (a *Context) AddParam(key string, value interface{}) {
	a.SetUserValue(key, value)
}

func (a *Context) ParamInt64(id string) int64 {
	return cast.Int64(a.UserValue(id))
}
func (a *Context) ParamInt(id string) int {
	return cast.Int(a.UserValue(id))
}

func (a *Context) ParamUint64(id string) uint64 {
	return cast.Uint64(a.UserValue(id))
}

func (a *Context) ParamUint32(id string) uint32 {
	return uint32(cast.Uint64(a.UserValue(id)))
}

func (a *Context) ParamBool(id string) bool {
	return cast.Bool(a.UserValue(id))
}
