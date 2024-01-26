package api

import (
	"net/http"
	"net/url"

	"github.com/TcMits/vnprovince/rest"
	"github.com/TcMits/vnprovince/rest/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/protobuf/encoding/protojson"
)

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

	proto.RegisterVNProvinceServiceHandlerServer(r.Context(), mux, rest.NewVNProvinceService())

	docsURL, _ := url.JoinPath("/api/static", rest.GetJSONDocsPath())
	mux.HandlePath(http.MethodGet, "/api/static/{filepath=**}", handlerFromHTTPHandler(http.StripPrefix("/api/static/", http.FileServer(rest.GetDocsFS()))))
	mux.HandlePath(http.MethodGet, "/api/docs", handlerFromHTTPHandler(rest.SwaggerDocs(docsURL)))

	var handler http.Handler = cors.AllowAll().Handler(mux)
	handler.ServeHTTP(w, r)
}
