package main

import (
	"encoding/json"
	"gopacs/templates"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/proxy.pac", proxyPac)
	log.SetFlags(log.Lshortfile)
}
func getGFWRole() []string {
	str, _ := readAllFileToByte(gfwlist_path)
	strs := strings.Split(string(str), "\n")
	return strs
}

type ProxyInfo struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Address  string `json:"address"`
}

type RoleList struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func proxyPac(w http.ResponseWriter, r *http.Request) {
	//t, err := template.ParseFiles("template/html/404.html")
	r.ParseForm()
	var err error
	gfw := getGFWRole()
	log.Println("a user connect this. ")

	var b []byte

	b, err = readAllFileToByte(proxy_server_path)
	if err != nil {
		http.Error(w, err.Error(), 502)
		log.Println("json...... ", err.Error())
		return
	}
	var v = []ProxyInfo{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		http.Error(w, err.Error(), 502)
		log.Println("proxy_error", err.Error())
		return
	}
	b, err = readAllFileToByte(role_list_path)
	if err != nil {
		http.Error(w, err.Error(), 502)
		log.Println("role_list_path ", err.Error())
		return
	}
	var role = []RoleList{}
	err = json.Unmarshal(b, &role)
	if err != nil {
		http.Error(w, err.Error(), 502)
		log.Println("role error", err.Error())

		return
	}
	ssh := []string{}

	getList := GetListPorts{make(chan []string)}
	go func() {
		select {
		case ssh_get_all_ports <- getList:
		default:
		}
	}()
	select {
	case ssh = <-getList.Info:
	case <-time.After(10 * time.Millisecond):
	}
	p := struct {
		Proxy  []ProxyInfo
		Role   []RoleList
		GFW    []string
		Server string
		Ssh    []string
	}{v, role, gfw, serverName, ssh}
	r.Form.Get("key")
	w.Header().Set("Server", "GoPacProxy 1.0")
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	err = templates.T.ExecuteTemplate(w, "proxy.pac.tpl", p)

	if err != nil {
		log.Println(err)
	}

}

func readAllFileToByte(name string) ([]byte, error) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0777)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	var b []byte
	b, err = ioutil.ReadAll(file)
	return b, err
}
