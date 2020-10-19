package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cernbox/cboxawareness/lbproxy"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"
)

var pushFlag bool
var carbonServerFlag string

func init() {
	flag.BoolVar(&pushFlag, "p", false, "push metrics to carbon metrics store")
	flag.StringVar(&carbonServerFlag, "carbon-server", "filer-carbon.cern.ch:2003", "carbon server")
	flag.Parse()
}

func main() {
	today := time.Now().Format("2006/01/02")

	// analyze lbproxy logs
	uniqUsersMetric := lbproxy.NewUniqUsersMetric()
	pattern := path.Join("/data/log/", today, "box.lbproxy/*/td.var.log.cboxredirectd.cboxredirectd_http.log")
	metrics := parse(pattern, uniqUsersMetric)

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
	format := "%s %d %d\n"
	for _, m := range metrics {
		for k, v := range m {
			payload := fmt.Sprintf(format, k, v, now)
			d(payload)
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
