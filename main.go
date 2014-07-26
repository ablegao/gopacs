package main

import (
	"flag"
	"gopacs/templates"
	"log"
	"net/http"
	"os/exec"
	"strings"
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

var serverName string

func main() {
	templates.Parse()
	flag.Parse()
	info, _ := exec.Command("uname", "-a").Output()
	str := strings.Split(string(info), " ")
	serverName = str[1]

	log.Printf("SERVER: http://%s%s\n", serverName, *address)
	log.Printf("PAC: http://%s%s/proxy.pac\n", serverName, *address)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Println(err)
	}
}
