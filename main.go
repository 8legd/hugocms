package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/8legd/hugocms/admin"

	"github.com/adrianduke/configr"
)

func main() {

	configr.AddSource(configr.NewFileSource("config/server.toml"))
	if err := configr.Parse(); err != nil {
		handleError(err)
	}
	port, err := configr.Int("port")
	if err != nil {
		handleError(err)
	}
	mux := http.NewServeMux()
	admin.Admin.MountTo("/admin", mux)
	// ./system is where QOR admin will upload files e.g. images
	for _, path := range []string{"css", "fonts", "images", "js", "system"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("public")))
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		handleError(err)
	}

	fmt.Printf("Listening on: %v\n", port)
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	//TODO more graceful exit!
}
