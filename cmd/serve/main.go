package main

import (
	"net/http"

	"github.com/TcMits/vnprovince/api"
)

func main() {
	if err := http.ListenAndServe(":8080", http.HandlerFunc(api.ServeHTTP)); err != nil {
		panic(err)
	}
}
