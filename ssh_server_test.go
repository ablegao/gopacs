package main

import (
	"testing"
	"time"
)

func init() {
	go ssh_controller()
}
func Test_SshStart(t *testing.T) {

	item := new(SshParams)

	item.Name = "myhost"
	item.Address = "root@50.93.203.167"
	item.ServerPort = "22"
	item.LocalPort = "8089"
	item.Passwd = "gaoenbo521"

	ssh_start_chan <- item

	<-time.After(40 * time.Second)
}
