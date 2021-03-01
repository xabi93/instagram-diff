package server

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"text/template"

	"github.com/xabi93/instagram-diff/instagram"
)

//go:embed index.tmpl
var tmplt []byte

func Serve(port string, differ *instagram.Instadiff) error {
	handler(differ)

	fmt.Printf("Showing result at: http://localhost:%s\n", port)

	return http.ListenAndServe(net.JoinHostPort("", port), nil)
}

func handler(d *instagram.Instadiff) {
	t, err := template.New("result").Parse(string(tmplt))
	if err != nil {
		log.Fatal(err)
	}
	type table struct {
		Title string
		Users []instagram.User
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result, err := d.Diff(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := t.Execute(w, []table{
			{Title: "Follow not follower", Users: result.FollowNotFollower},
			{Title: "Follower not follow", Users: result.FollowerNotFollow},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
}
