package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"text/template"

	instadiff "github.com/xabi93/instagram-diff"
)

func Serve(port string, differ *instadiff.Instadiff) error {
	handler(differ)

	fmt.Printf("Serving, at port: %s\n", port)

	return http.ListenAndServe(net.JoinHostPort("", port), nil)
}

func handler(d *instadiff.Instadiff) {
	t, err := template.New("result").Parse(resultTmpl)
	if err != nil {
		log.Fatal(err)
	}
	type table struct {
		Title string
		Users []instadiff.User
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result, err := d.Diff(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := t.Execute(w, []table{
			table{Title: "Follow not follower", Users: result.FollowNotFollower},
			table{Title: "Follower not follow", Users: result.FollowerNotFollow},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
}

const resultTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Follow/Followers</title>
		<style>
		body{
			background:#eee;
		}
		.main-box.no-header {
			padding-top: 20px;
		}
		.main-box {
			background: #FFFFFF;
			-webkit-box-shadow: 1px 1px 2px 0 #CCCCCC;
			-moz-box-shadow: 1px 1px 2px 0 #CCCCCC;
			-o-box-shadow: 1px 1px 2px 0 #CCCCCC;
			-ms-box-shadow: 1px 1px 2px 0 #CCCCCC;
			box-shadow: 1px 1px 2px 0 #CCCCCC;
			margin-bottom: 16px;
			-webikt-border-radius: 3px;
			-moz-border-radius: 3px;
			border-radius: 3px;
		}
		.main-box-body {
			display: flex;
			justify-content: space-evenly;
		}
		.table a.table-link.danger {
			color: #e74c3c;
		}
		.label {
			border-radius: 3px;
			font-size: 0.875em;
			font-weight: 600;
		}
		.user-list tbody td .user-subhead {
			font-size: 0.875em;
			font-style: italic;
		}
		.user-list tbody td .user-link {
			display: block;
			font-size: 1.25em;
			padding-top: 3px;
			margin-left: 60px;
		}
		a {
			color: #3498db;
			outline: none!important;
		}
		.user-list tbody td>img {
			position: relative;
			max-width: 50px;
			float: left;
			margin-right: 15px;
		}

		.table thead tr th {
			text-transform: uppercase;
			font-size: 0.875em;
		}
		.table thead tr th {
			border-bottom: 2px solid #e7ebee;
		}
		.table tbody tr td:first-child {
			font-size: 1.125em;
			font-weight: 300;
		}
		.table tbody tr td {
			font-size: 0.875em;
			vertical-align: middle;
			border-top: 1px solid #e7ebee;
			padding: 12px 8px;
		}
		</style>
	</head>
	<body>
		<div class="container bootstrap snippet">
			<div class="row">
				<div class="col-lg-12">
					<div class="main-box no-header clearfix">
						<div class="main-box-body clearfix">
							{{range .}}
								{{template "userList" .}}
							{{end}}
						</div>
					</div>
				</div>
			</div>
		</div>
	</body>
</html>

{{define "userList"}}
	<div class="table-responsive">
	<h2>{{.Title}}</h2>
	<table class="table user-list">
		<tbody>
		{{range .Users}}
			<tr>
				<td>
					<img src="{{.ProfilePic}}">
					<a href="https://www.instagram.com/{{.Username}}" target="_blank" class="user-link">{{.Username}}</a>
					<span class="user-subhead">{{.Fullname}}</span>
				</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	</div>
{{end}}
`
