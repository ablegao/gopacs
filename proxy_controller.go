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

}

func proxy_ctrl(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("template/html/404.html")
	w.Header().Set("Server", "GoPacProxy 1.0")
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case "GET":
		file, err := os.OpenFile(proxy_server_path, os.O_RDONLY|os.O_CREATE, 0777)
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
		file, err := os.OpenFile(proxy_server_path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
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

func params_ctrl(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("template/html/404.html")
	w.Header().Set("Server", "GoPacProxy 1.0")
	w.Header().Set("Content-Type", "application/json")
	filePath := role_list_path
	switch r.Method {

	case "GET":
		file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0777)
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
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer file.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		file.Write(b)
		w.Write([]byte(""))

	}

}
