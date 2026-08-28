package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var distEmbed embed.FS

// Handler 返回内嵌前端静态资源；未知路径回退 index.html（SPA）。
func Handler() (http.Handler, error) {
	sub, err := fs.Sub(distEmbed, "dist")
	if err != nil {
		return nil, err
	}
	return spa(sub), nil
}

func spa(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if strings.Contains(path, "..") {
			http.NotFound(w, r)
			return
		}
		// 带扩展名的资源：存在则直出，不存在 404（避免把 .js 误回退成 HTML）
		if strings.Contains(path, ".") {
			if _, err := fs.Stat(fsys, path); err != nil {
				http.NotFound(w, r)
				return
			}
			fileServer.ServeHTTP(w, r)
			return
		}
		// 路由页：无实体文件则回退 index.html
		if _, err := fs.Stat(fsys, path); err != nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/"
			fileServer.ServeHTTP(w, r2)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
