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

type SyncDistrMetric struct {
	// cernbox.awareness.sync.distr.cbox.2-7-0.windows 123
	dist map[string]int
}

func (uu *SyncDistrMetric) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
		return
	}

	agent := strings.ToLower(l.Agent)
	if !strings.Contains(agent, "mirall") {
		return
	}

	// most lines to parse are like these:
	// 0           1         2            3      4      5
	// mozilla/5.0 (windows) mirall/2.4.2 (build 1396) (cernbox)
	// but sometimes line can avoid build info:
	// 0           1       2            3
	// mozilla/5.0 (linux) mirall/2.6.3 (cernbox)
	// mozilla/5.0 (linux) mirall/2.6.4git (nextcloud)
	tokens := strings.Split(agent, " ")

	os := tokens[1]
	os = strings.TrimPrefix(os, "(")
	os = strings.TrimSuffix(os, ")")

	version := strings.ReplaceAll(strings.Split(tokens[2], "/")[1], ".", "-")

	var brand string
	if len(tokens) == 3 {
		brand = "owncloud"
	} else {
		b := strings.TrimPrefix(tokens[3], "(")
		if strings.HasPrefix(b, "build") {
			// check for empty platform, we map to owncloud
			if len(tokens) < 6 {
				b = "owncloud"
			} else {
				b = strings.TrimSuffix(strings.TrimPrefix(tokens[5], "("), ")")
			}
		} else {
			b = strings.TrimSuffix(b, ")")
		}

		if b == "" {
			panic("b is empty")
		}
		brand = b
	}

	uu.dist[fmt.Sprintf("cernbox.awareness.sync.dist.%s.%s.%s", os, version, brand)]++
}

func NewSyncDistrMetric() *SyncDistrMetric {
	return &SyncDistrMetric{
		dist: map[string]int{},
	}
}

func (uu *SyncDistrMetric) Metrics() map[string]int {
	return uu.dist
}

func er(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v", err)
}
