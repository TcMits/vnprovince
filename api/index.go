package api

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/TcMits/vnprovince/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed proto/*.json
var embedAssets embed.FS

func swaggerDocs(assetPath string) runtime.HandlerFunc {
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
        url: '/api/static/{{.}}',
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
	t.Execute(&b, assetPath)
	content := b.String()
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	}
}

func handlerFromHTTPHandler(h http.Handler) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		h.ServeHTTP(w, r)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	proto.RegisterVNProvinceServiceHandlerServer(r.Context(), mux, newVNProvinceService())
	mux.HandlePath(http.MethodGet, "/api/static/{filepath=**}", handlerFromHTTPHandler(http.StripPrefix("/api/static/", http.FileServer(http.FS(embedAssets)))))
	mux.HandlePath(http.MethodGet, "/api/docs", swaggerDocs("proto/docs.swagger.json"))

	var handler http.Handler = cors.AllowAll().Handler(mux)
	handler.ServeHTTP(w, r)
}
