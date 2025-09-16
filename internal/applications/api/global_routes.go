package api

import (
	"net/http"

	_ "github.com/edalferes/monogo/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterGlobalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("metrics: not implemented"))
	})

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

}
