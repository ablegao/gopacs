package main

import (
	"flag"
	"gopacs/templates"
	"log"
	"net/http"
)

var (
	db_path           = "./"
	proxy_server_path = "./proxy.json"
	role_list_path    = "./params.json"
	gfwlist_path      = "./gfwlist.role"
)

var address = flag.String("address", ":8888", "0.0.0.0:8888")

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", home)

}
func home(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("template/html/404.html")

	err := templates.T.ExecuteTemplate(w, "home.tpl", nil)

	if err != nil {
		log.Println(err)
	}

}

func main() {
	templates.Parse()
	flag.Parse()

	//http.HandleFunc("/admin/", adminHandler)
	//http.HandleFunc("/login/", loginHandler)
	//http.HandleFunc("/ajax/", ajaxHandler)
	//http.HandleFunc("/", NotFoundHandler)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Println(err)
	}
}
