package cmd

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"tiny-gateway/api"
	"tiny-gateway/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var RootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "tiny gateway",
	Long: "tiny gateway",
	Run: run,
}


func run(cmd *cobra.Command, args[] string) {
	apiDefinitions, err := api.LoadDefinitions(config.Load().ConfigDir)
	if err != nil {
		log.Panic(err)
	}

	router := chi.NewRouter()

	for _, definition := range apiDefinitions.GetDefinitions(){
		methods := definition.Proxy.Methods

		director := func(req *http.Request) {
			target, _ := url.Parse(definition.Proxy.UpstreamURL)
			//targetQuery := target.RawQuery

			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			path := target.Path
			req.URL.Path = path

			req.Host = target.Host
			req.RequestURI = ""

		}

		reverseHttp := &httputil.ReverseProxy{
			Director:       director,
			Transport:      http.DefaultTransport,
		}

		for _, method := range methods {
			if method == "ANY" {
				router.Handle(definition.Proxy.ListenPath, reverseHttp)
			} else {
				router.Method(method, definition.Proxy.ListenPath, reverseHttp)
			}
		}
	}

	if err = http.ListenAndServe(":1997", router); err != nil {
		panic(err)
	}
}
