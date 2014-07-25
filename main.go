package main

import (
	"gopacs/templates"
	"log"
	"net/http"
)

var (
	proxy_server_path = "./proxy.json"
	role_list_path    = "./params.json"
	gfwlist_path      = "./gfwlist.role"
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", home)

}

func main() {
	templates.Parse()

	//http.HandleFunc("/admin/", adminHandler)
	//http.HandleFunc("/login/", loginHandler)
	//http.HandleFunc("/ajax/", ajaxHandler)
	//http.HandleFunc("/", NotFoundHandler)

	http.ListenAndServe(":8888", nil)

}

func home(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("template/html/404.html")

	err := templates.T.ExecuteTemplate(w, "home.tpl", nil)

	if err != nil {
		log.Println(err)
	}

}
