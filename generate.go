//go:build ignore

package vnprovince

//go:generate go run ./cmd/generate

//go:generate mkdir -p v2
//go:generate cp go.mod v2/go.mod
//go:generate go mod edit -module github.com/TcMits/vnprovince/v2 v2/go.mod
//go:generate go run ./cmd/generatev2
