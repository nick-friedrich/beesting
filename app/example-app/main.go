package main

import (
	"embed"
	"net/http"

	"github.com/nick-friedrich/beesting/pkg/beesting"
)

//go:embed templates static
var fs embed.FS

func main() {
	app := beesting.NewApp(fs)
	app.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from BeeSting app!"))
	})
	app.Run(":8080")
}
