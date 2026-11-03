// Package nethttp 是标准库 net/http 的 GoAdmin 适配器：路由挂在 *http.ServeMux 上（Go 1.22 的 "GET /path/{name}" 模式）。
// 裁剪时把 gin、chi 等十几个框架的适配器都删了，只留这一个；写法照着原来的 chi 适配器。
//
//	mux := http.NewServeMux()
//	eng := engine.Default()
//	err := eng.AddConfig(&cfg).AddGenerators(...).Use(mux)
package nethttp

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"goadmin/adapter"
	"goadmin/context"
	"goadmin/engine"
	cfg "goadmin/modules/config"
	"goadmin/plugins"
	"goadmin/plugins/admin/models"
	"goadmin/plugins/admin/modules/constant"
	"goadmin/template/types"
)

// NetHTTP 是 net/http 的 GoAdmin 适配器。
type NetHTTP struct {
	adapter.BaseAdapter
	ctx Context
	app *http.ServeMux
}

func init() {
	engine.Register(new(NetHTTP))
}

// User implements the method Adapter.User.
func (n *NetHTTP) User(ctx interface{}) (models.UserModel, bool) {
	return n.GetUser(ctx, n)
}

// Use implements the method Adapter.Use.
func (n *NetHTTP) Use(app interface{}, plugs []plugins.Plugin) error {
	return n.GetUse(app, plugs, n)
}

// Content implements the method Adapter.Content.
func (n *NetHTTP) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	n.GetContent(ctx, getPanelFn, n, btns, fn)
}

// SetApp implements the method Adapter.SetApp.
func (n *NetHTTP) SetApp(app interface{}) error {
	mux, ok := app.(*http.ServeMux)
	if !ok {
		return errors.New("nethttp adapter SetApp: wrong parameter, want *http.ServeMux")
	}
	n.app = mux
	return nil
}

// paramPattern 是 GoAdmin 路由里的 :name 参数。
var paramPattern = regexp.MustCompile(`:([^/]+)`)

// AddHandler implements the method Adapter.AddHandler：
// GoAdmin 的 "/info/:__prefix" 转成 ServeMux 的 "GET /info/{__prefix}"，路径参数照原来的适配器追加到查询串里（GoAdmin 从查询串读）。
func (n *NetHTTP) AddHandler(method, path string, handlers context.Handlers) {
	if len(path) > 1 && path[0] == '/' && path[1] == '/' {
		path = path[1:]
	}
	var names []string
	for _, m := range paramPattern.FindAllStringSubmatch(path, -1) {
		names = append(names, m[1])
	}
	pattern := paramPattern.ReplaceAllString(path, "{$1}")
	if strings.HasSuffix(pattern, "/") {
		pattern += "{$}" // 以 / 结尾的模式在 ServeMux 里匹配整棵子树，这里只要这一个路径
	}
	if method = strings.ToUpper(method); method != "ANY" {
		pattern = method + " " + pattern
	}

	n.app.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}
		if len(names) > 0 {
			query := r.URL.Query()
			for _, name := range names {
				query.Add(name, r.PathValue(name))
			}
			r.URL.RawQuery = query.Encode()
		}

		ctx := context.NewContext(r)
		ctx.SetHandlers(handlers).Next()
		for key, head := range ctx.Response.Header {
			w.Header().Set(key, head[0])
		}
		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			w.WriteHeader(ctx.Response.StatusCode)
			_, _ = w.Write(buf.Bytes())
		} else {
			w.WriteHeader(ctx.Response.StatusCode)
		}
	})
}

// Context 是一次请求的 Request 和 ResponseWriter。
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// SetContext implements the method Adapter.SetContext.
func (*NetHTTP) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	ctx, ok := contextInterface.(Context)
	if !ok {
		panic("nethttp adapter SetContext: wrong parameter")
	}
	return &NetHTTP{ctx: ctx}
}

// Name implements the method Adapter.Name.
func (*NetHTTP) Name() string {
	return "nethttp"
}

// Redirect implements the method Adapter.Redirect.
func (n *NetHTTP) Redirect() {
	http.Redirect(n.ctx.Response, n.ctx.Request, cfg.Url(cfg.GetLoginUrl()), http.StatusFound)
}

// SetContentType implements the method Adapter.SetContentType.
func (n *NetHTTP) SetContentType() {
	n.ctx.Response.Header().Set("Content-Type", n.HTMLContentType())
}

// Write implements the method Adapter.Write.
func (n *NetHTTP) Write(body []byte) {
	n.ctx.Response.WriteHeader(http.StatusOK)
	_, _ = n.ctx.Response.Write(body)
}

// GetCookie implements the method Adapter.GetCookie.
func (n *NetHTTP) GetCookie() (string, error) {
	cookie, err := n.ctx.Request.Cookie(n.CookieKey())
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

// Lang implements the method Adapter.Lang.
func (n *NetHTTP) Lang() string {
	return n.ctx.Request.URL.Query().Get("__ga_lang")
}

// Path implements the method Adapter.Path.
func (n *NetHTTP) Path() string {
	return n.ctx.Request.URL.Path
}

// Method implements the method Adapter.Method.
func (n *NetHTTP) Method() string {
	return n.ctx.Request.Method
}

// FormParam implements the method Adapter.FormParam.
func (n *NetHTTP) FormParam() url.Values {
	_ = n.ctx.Request.ParseMultipartForm(32 << 20)
	return n.ctx.Request.PostForm
}

// IsPjax implements the method Adapter.IsPjax.
func (n *NetHTTP) IsPjax() bool {
	return n.ctx.Request.Header.Get(constant.PjaxHeader) == "true"
}

// Query implements the method Adapter.Query.
func (n *NetHTTP) Query() url.Values {
	return n.ctx.Request.URL.Query()
}

// Request implements the method Adapter.Request.
func (n *NetHTTP) Request() *http.Request {
	return n.ctx.Request
}
