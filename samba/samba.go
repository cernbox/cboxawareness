package samba

import (
	"fmt"
	"github.com/json-iterator/go"
	"os"
	"strings"
	"time"
)

type line struct {
	Host        string    `json:"host"`
	Ident       string    `json:"ident"`
	Pid         string    `json:"pid"`
	Message     string    `json:"message"`
	File        string    `json:"file"`
	Tag         string    `json:"tag"`
	Hostname    string    `json:"hostname"`
	Hostgroup   string    `json:"hostgroup"`
	Shostgroup  string    `json:"shostgroup"`
	Environment string    `json:"environment"`
	Time        time.Time `json:"time"`
}

type UniqUsers struct {
	dist map[string]int
}

func (uu *UniqUsers) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}

	if !strings.Contains(l.Message, "opened file") && !strings.Contains(l.Message, "closed file") {
		return
	}

	tokens := strings.SplitN(l.Message, " ", 2)
	username := tokens[0]
	uu.dist[username]++
}

func NewUniqUsers() *UniqUsers {
	return &UniqUsers{dist: map[string]int{}}
}

func (uu *UniqUsers) Metrics() map[string]int {
	m := map[string]int{
		"users.samba": len(uu.dist),
	}
	return m
}

func er(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v", err)
}
