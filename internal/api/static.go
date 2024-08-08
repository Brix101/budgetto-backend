package api

import (
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"go.uber.org/zap"
)

func getStaticHandler(fSys fs.FS, logger *zap.Logger) (index, static http.HandlerFunc) {

	index = func(w http.ResponseWriter, r *http.Request) {

		// Check if the request path contains "/api/"
		if strings.Contains(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		data := domain.StaticData{
			Name: "YourNameHere",
		}

		fileContents, err := fs.ReadFile(fSys, "index.html")
		if err != nil {
			logger.Error("failed to read file 'index.html' from assets filesystem", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.New("index").Delims("[{[", "]}]").Parse(string(fileContents)))
		if err := tmpl.Execute(w, data); err != nil {
			logger.Error("failed to execute template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	static = http.FileServer(http.FS(fSys)).ServeHTTP

	return index, static
}
