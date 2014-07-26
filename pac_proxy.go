package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopacs/templates"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var coder = base64.StdEncoding

func init() {
	http.HandleFunc("/proxy.pac", proxyPac)
	log.SetFlags(log.Lshortfile)
}
func getGFWRole() []string {
	str, _ := readAllFileToByte(gfwlist_path)
	strs := strings.Split(string(str), "\n")
	return strs
}

func getGFWList() ([]string, error) {
	if b, err := readAllFileToByte(gfwlist_path); err != nil {
		return []string{}, err
	} else {
		str, err := coder.DecodeString(string(b))
		if err != nil {
			return []string{}, err
		}
		strs := strings.Split(string(str), "\n")
		for _, str := range strs {
			str = strings.Trim(str, "\t \n")
			if len(str) == 0 {
				continue
			}
			if str[0] == '!' {
				continue
			}
			if str[0] == '@' {
				continue
			}
			str = strings.Replace(str, ".", "\\.", -1)
			str = strings.Replace(str, "/", "\\/", -1)
			str = strings.Replace(str, "%", "\\%", -1)
			str = strings.Replace(str, ":", "\\:", -1)

			if string(str[0:2]) == "||" {
				str = strings.Replace(str, "||", "^[\\w\\-]+:\\/+(?!\\/)(?:[^\\/]+\\.)?", -1)
			}

			if string(str[0]) == "|" {
				str = strings.Replace(str, "|", "^", 1)
			}
			str = "/" + str + "/i"
			fmt.Println(str)

		}
		return strs, nil
	}

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
	p := struct {
		Proxy []ProxyInfo
		Role  []RoleList
		GFW   []string
	}{v, role, gfw}

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
