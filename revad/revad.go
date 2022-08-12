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
	Path        string    `json:"path"`
	Agent       string    `json:"agent"`
	User        string    `json:"user"`
	Country     string    `json:"country"`
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

//Apps used metric scanning for OCIS specific app identifiers
type OCISAppsMetric struct {
	dist map[string]int
}

func (uu *OCISAppsMetric) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}

	if strings.Contains(l.Path, "/text-editor") {
		uu.dist["OCIS.apps.usage.OCIStext-editor"]++
	} else if strings.Contains(l.Path, "/*.docx?app=MS'%'20365'%'20on'%'20Cloud") {
		uu.dist["OCIS.apps.usage.word"]++
	} else if strings.Contains(l.Path, "/draw-io") {
		uu.dist["OCIS.apps.usage.OCISdrawio"]++
	} else if strings.Contains(l.Path, "/*.md?app=CodiMD") {
		uu.dist["OCIS.apps.usage.CodiMD"]++
	} else if strings.Contains(l.Path, "/*.xlsx?app=MS'%'20365'%'20on'%'20Cloud") {
		uu.dist["OCIS.apps.usage.excel"]++
	} else if strings.Contains(l.Path, "/*.pptx?app=MS'%'20365'%'20on'%'20Cloud") {
		uu.dist["OCIS.apps.usage.pwpoint"]++
	} else if strings.Contains(l.Path, "/*.odt?app=Collabora") {
		uu.dist["OCIS.apps.usage.openDoc"]++
	} else if strings.Contains(l.Path, "/*.ods?app=Collabora") {
		uu.dist["OCIS.apps.usage.openSpread"]++
	} else if strings.Contains(l.Path, "/*.odp?app=Collabora") {
		uu.dist["OCIS.apps.usage.openSlide"]++
	}

}

func NewOCISAppsMetric() *OCISAppsMetric {
	return &OCISAppsMetric{
		dist: map[string]int{},
	}
}

func (uu *OCISAppsMetric) Metrics() map[string]int {
	return uu.dist
}

//metric that keeps track of Unique OCIS Users 
type OCISUniqUsers struct {
	dist map[string]int
}

func (uu *OCISUniqUsers) Do(data []byte) {
	l := &line{}

	if err := jsoniter.Unmarshal([]byte(data), l); err != nil {
		er(err)
	}

	if l.User != "" {
		uu.dist[l.User]++
	}
}

func NewOCISUniqUsers() *OCISUniqUsers {
	return &OCISUniqUsers{dist: map[string]int{}}
}

func (uu *OCISUniqUsers) Metrics() map[string]int {
	m := map[string]int{
		"users.OCISweb": len(uu.dist),
	}
	return m
}