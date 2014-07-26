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
	ssh_start_chan  = make(chan *SshParams)
	ssh_stop_chan   = make(chan string)
	ssh_run_server  = map[string]*SshParams{}
	ssh_stop_server = make(chan bool)
)

func init() {
	http.HandleFunc("/ssh_start", ssh_start)
	http.HandleFunc("/ssh_stop", ssh_stop)
	go ssh_controller()
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
		}
		s.State = 0

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
			err := cmd.Wait()
			log.Println(err)
			s.stop <- true
		}()

		for {
			select {
			case <-s.stop:

				return
			case <-s.kill: //终止服务
				if cmd.ProcessState.Success() {
					cmd.Process.Kill()
				}
				s.State = 4
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
				log.Println("服务开启调通道被关闭")
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
				log.Println("服务终止通道被关闭")
				return
			}
			_, ok = ssh_run_server[name]
			if !ok {
				log.Println("服务没有启动过")
				return
			}
			if ssh_run_server[name].State == 2 {
				ssh_run_server[name].kill <- true
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

}
