package lbproxy

import (
	"fmt"
	"github.com/json-iterator/go"
	"os"
	"strings"
	"time"
)

type line struct {
	Host        string    `json:"host"`
	User        string    `json:"user"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Code        int       `json:"code"`
	Size        int       `json:"size"`
	Referer     string    `json:"referer"`
	Agent       string    `json:"agent"`
	File        string    `json:"file"`
	Tag         string    `json:"tag"`
	Hostname    string    `json:"hostname"`
	Hostgroup   string    `json:"hostgroup"`
	Shostgroup  string    `json:"shostgroup"`
	Environment string    `json:"environment"`
	Time        time.Time `json:"time"`
}

type UniqUsersMetric struct {
	desktop map[string]int
	mobile  map[string]int
}

func (uu *UniqUsersMetric) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}
	// filter out public links and empty ones
	if l.User == "" || strings.HasPrefix(l.Path, "/public.php/webdav/") {
		return
	}

	// classify based on desktop, mobile and web access
	if strings.Contains(l.Path, "cernbox/desktop") {
		uu.desktop[l.User]++
	} else if strings.Contains(l.Path, "cernbox/mobile") {
		uu.mobile[l.User]++
	}
}

func NewUniqUsersMetric() *UniqUsersMetric {
	return &UniqUsersMetric{
		desktop: map[string]int{},
		mobile:  map[string]int{},
	}
}

func (uu *UniqUsersMetric) Metrics() map[string]int {
	m := map[string]int{}
	m["cernbox.awareness.users.desktop"] = len(uu.desktop)
	m["cernbox.awareness.users.mobile"] = len(uu.mobile)
	return m
}

func er(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v", err)
}
