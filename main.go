package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cernbox/cboxawareness/lbproxy"
	"github.com/cernbox/cboxawareness/revad"
	"github.com/cernbox/cboxawareness/samba"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var pushFlag bool
var prefixFlag string
var carbonServerFlag string
var pickFlag string

func init() {
	flag.BoolVar(&pushFlag, "p", false, "push metrics to carbon metrics store")
	flag.StringVar(&carbonServerFlag, "carbon-server", "filer-carbon.cern.ch:2003", "carbon server")
	flag.StringVar(&prefixFlag, "prefix", "cernbox.awareness", "carbon metrics prefix")
	flag.StringVar(&pickFlag, "pick", "", "pick only one file to parse")
	flag.StringVar(&pathFlag, "path", "/data/cephfs/logs/", "path where the logs are stored")
	flag.Parse()
}

func main() {
	today := time.Now().Format("2006/01/02")

	// analyze lbproxy logs
	uniqUsersMetric := lbproxy.NewUniqUsersMetric()
	syncDistrMetric := lbproxy.NewSyncDistrMetric()
	countryMetric := lbproxy.NewCountryMetric()
	appsMetrics := lbproxy.NewAppsMetric()
	pattern := path.Join(pathFlag, today, "box.lbproxy/*/td.var.log.cboxredirectd.cboxredirectd_http.log")
	metrics := parse(pattern, uniqUsersMetric, syncDistrMetric, countryMetric, appsMetrics)
	// analyze revad logs
	userCreated := revad.NewUserCreated()
	uniqWeb := revad.NewUniqUsers()
	pattern = path.Join(pathFlag, today, "box.*/*/td.var.log.revad.revad_app.log")
	metrics = append(metrics, parse(pattern, userCreated, uniqWeb)...)

	// analyze samba logs
	uniqSamba := samba.NewUniqUsers()
	pattern = path.Join(pathFlag, today, "box.samba*/*/td.var.log.samba.smbclients.log")
	metrics = append(metrics, parse(pattern, uniqSamba)...)

	// send to carbon
	push(metrics)
}

func push(metrics []map[string]int) {
	// create tcp connection
	conn, err := net.Dial("tcp", carbonServerFlag)
	if err != nil {
		er(err)
	}

	now := time.Now().Unix()
	format := "%s.%s %d %d\n"
	for _, m := range metrics {
		for k, v := range m {
			payload := fmt.Sprintf(format, prefixFlag, k, v, now)
			fmt.Print(payload)
			if _, err := conn.Write([]byte(payload)); err != nil {
				er(err)
			}
		}
	}
}

func parse(pattern string, metrics ...Metric) (counters []map[string]int) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		er(err)
	}

	for _, f := range files {
		if pickFlag != "" && !strings.Contains(f, pickFlag) {
			continue
		}
		fd, err := os.Open(f)
		if err != nil {
			er(err)
		}
		defer fd.Close()

		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			data := scanner.Text()

			for _, m := range metrics {
				m.Do([]byte(data))
			}
		}

		if err := scanner.Err(); err != nil {
			er(err)
		}

	}

	for _, m := range metrics {
		counters = append(counters, m.Metrics())
	}

	return
}

type Metric interface {
	Do(data []byte)
	Metrics() map[string]int
}

func er(err error) {
	fmt.Fprintf(os.Stderr, "error: %+v", err)
}

func d(v interface{}) {
	fmt.Println(v)
}
