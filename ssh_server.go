package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var (
	ssh_start_chan    = make(chan *SshParams)
	ssh_stop_chan     = make(chan string)
	ssh_run_server    = map[string]*SshParams{}
	ssh_stop_server   = make(chan bool)
	ssh_get_status    = make(chan GetListStatus)
	ssh_get_all_ports = make(chan GetListPorts)
)

func init() {
	http.HandleFunc("/ssh_start", ssh_start)
	http.HandleFunc("/ssh_stop", ssh_stop)
	http.HandleFunc("/ssh_state", ssh_status)

	go ssh_controller()
}

type GetListStatus struct {
	Info chan map[string]int
}

type GetListPorts struct {
	Info chan []string
}

type SshParams struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	ServerPort string `json:"server_port"`
	LocalPort  string `json:"local_port"`
	Passwd     string `json:"passwd"`
	State      int    `json:"state"`

	stop        chan bool
	kill        chan bool
	changeState chan int
	restartNum  int
}

func (s *SshParams) Start() {
	t := time.NewTicker(1. * time.Second)
	s.State = 2
	defer func() {
		if s.State != 4 && s.restartNum < 2 {
			go s.Start()
			s.restartNum = s.restartNum + 1
			log.Println("Start error , try agent. ", s.restartNum)
			return
		}
		s.State = 0
		log.Println("Server can't run!")
	}()

	//strcmd = "/usr/bin/ssh"

	cmd := exec.Command("/usr/bin/ssh", "-C", "-o",
		"ServerAliveInterval=15",
		"-o",
		"ServerAliveCountMax=30",
		"-N",
		"-D",
		fmt.Sprintf("*:%v", s.LocalPort),
		s.Address,
		"-p", s.ServerPort,
		"'vmstat 30'",
	)

	log.Println(cmd.Args)

	if err := cmd.Start(); err == nil {

		go func() {
			s.State = 2
			err := cmd.Wait()
			log.Println(err)
			select {
			case s.stop <- true:
			default:
				log.Println("ssh is stop .bug not send s.stop <- true")
			}
		}()

		for {
			select {
			case <-s.stop:

				return
			case <-s.kill: //终止服务
				s.State = 4
				cmd.Process.Kill()
				return
			case id := <-s.changeState:
				s.State = id

			case <-t.C:

			}
		}

	} else {
		log.Println(err)

	}

}

func ssh_controller() {
	for {
		select {
		case <-ssh_stop_server:
			for _, item := range ssh_run_server {
				select {
				case item.kill <- true:
				default:
				}
			}
		case item, ok := <-ssh_start_chan:
			if !ok {
				log.Println("Channel is close ")
				return
			}
			var obj *SshParams
			obj, ok = ssh_run_server[item.Name]

			if !ok {
				ssh_run_server[item.Name] = item
				item.changeState = make(chan int)
				item.stop = make(chan bool)
				item.kill = make(chan bool)
				obj = item
			}
			if obj.State == 0 {

				go obj.Start()

				obj.changeState <- 2

			}
		case name, ok := <-ssh_stop_chan:
			if !ok {
				log.Println("Channel is close ")
				return
			}
			_, ok = ssh_run_server[name]
			if !ok {
				log.Println("Server into not exists! ")
				return
			}
			if ssh_run_server[name].State == 2 {
				ssh_run_server[name].kill <- true
			}
		case status := <-ssh_get_status:
			var m = map[string]int{}
			for name, item := range ssh_run_server {
				m[name] = item.State
			}
			select {
			case status.Info <- m:
			default:
			}
		case status := <-ssh_get_all_ports:
			var m = []string{}
			for _, item := range ssh_run_server {
				m = append(m, item.LocalPort)
			}
			select {
			case status.Info <- m:
			default:
			}
		}
	}

}

func ssh_start(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	defer r.Body.Close()
	v := SshParams{}

	err = json.Unmarshal(b, &v)

	if err != nil {
		http.Error(w, err.Error()+" "+string(b), 502)
		return
	}

	select {
	case ssh_start_chan <- &v:
		log.Println("start ssh ", v.Name)
	case <-time.After(1 * time.Second):
		log.Println("start ssh timeout  ", v.Name)

	}

}

func ssh_stop(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	defer r.Body.Close()
	select {
	case ssh_stop_chan <- string(b):
	default:

	}
}

func ssh_status(w http.ResponseWriter, r *http.Request) {
	status := GetListStatus{make(chan map[string]int)}
	go func() {
		select {
		case ssh_get_status <- status:
		default:
		}
	}()

	select {
	case out := <-status.Info:

		b, err := json.Marshal(out)
		if err != nil {
			http.Error(w, err.Error(), 502)
		}
		w.Write(b)
	case <-time.After(3):
	}
}
