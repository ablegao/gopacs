package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/proxy_ctrl", proxy_ctrl)
	http.HandleFunc("/params_ctrl", params_ctrl)
	http.HandleFunc("/ssh_ctrl", ssh_ctrl)

}

func proxy_ctrl(w http.ResponseWriter, r *http.Request) {

	dbapi(w, r, proxy_server_path)

}

func params_ctrl(w http.ResponseWriter, r *http.Request) {
	dbapi(w, r, role_list_path)

}
func ssh_ctrl(w http.ResponseWriter, r *http.Request) {

	dbapi(w, r, db_path+"ssh.json")
}
func dbapi(w http.ResponseWriter, r *http.Request, file_path string) {
	//t, err := template.ParseFiles("template/html/404.html")
	w.Header().Set("Server", "GoPacProxy 1.0")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	r.ParseForm()
	switch r.Method {

	case "GET":
		file, err := os.OpenFile(file_path, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			//w.Write([]byte(fmt.Sprintf("{'%s'}", err.Error())))
			return
		}
		defer file.Close()
		b, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.Write(b)
	case "POST":
		log.Println("执行了一次存储")
		file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			log.Println("open file error ", err.Error())

			return
		}
		defer file.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			log.Println("save error ", err.Error())
			return
		}
		file.Write(b)
		w.Write([]byte(""))

	}

}
