package rest

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed proto/*.json
var embedAssets embed.FS

func SwaggerDocs(docsURL string) http.HandlerFunc {
	t, _ := template.New("swagger").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta
      name="description"
      content="SwaggerUI"
    />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '{{.}}',
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout",
      });
    };
  </script>
  </body>
</html>
`)

	var b strings.Builder
	t.Execute(&b, docsURL)
	content := b.String()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	}
}

func GetDocsFS() http.FileSystem {
	return http.FS(embedAssets)
}

func GetJSONDocsPath() string {
	return "proto/docs.swagger.json"
}
