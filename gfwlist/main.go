package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var coder = base64.StdEncoding

func main() {
	getGFWList()
}
func getGFWList() ([]string, error) {
	if b, err := readAllFileToByte("./gfwlist.txt"); err != nil {
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
			if str[0] == '|' {
				continue
			}
			if str[0] == '.' {
				continue
			}
			/*
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
				str = "/" + str + "/i" */
			fmt.Println(str)

		}
		return strs, nil
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
