package revad

import (
	"fmt"
	"github.com/json-iterator/go"
	"os"
	"strings"
	"time"
)

type line struct {
	Level       string    `json:"level"`
	Caller      string    `json:"caller"`
	Msg         string    `json:"msg"`
	Args        string    `json:"args"`
	ExistStatus int       `json:"exist_status"`
	URL         string    `json:"url"`
	File        string    `json:"file"`
	Tag         string    `json:"tag"`
	Hostname    string    `json:"hostname"`
	Hostgroup   string    `json:"hostgroup"`
	Shostgroup  string    `json:"shostgroup"`
	Environment string    `json:"environment"`
	Time        time.Time `json:"time"`
	Username    string    `json:"username"`
}

type UserCreated struct {
	count int
}

func (uu *UserCreated) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}

	if strings.Contains(l.Msg, "homedir create for user") {
		uu.count++
	}
}

func NewUserCreated() *UserCreated {
	return &UserCreated{}
}

func (uu *UserCreated) Metrics() map[string]int {
	m := map[string]int{
		"newusers": uu.count,
	}
	return m
}

func er(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v", err)
}

type UniqUsers struct {
	dist map[string]int
}

func (uu *UniqUsers) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}

	if l.Username != "" {
		uu.dist[l.Username]++
	}
}

func NewUniqUsers() *UniqUsers {
	return &UniqUsers{dist: map[string]int{}}
}

func (uu *UniqUsers) Metrics() map[string]int {
	m := map[string]int{
		"users.web": len(uu.dist),
	}
	return m
}
