package main

import (
	"fmt"
	"io"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Stat struct {
	Name          string
	ReqCount      int
	RespCount     int
	RespNetErr    int
	RespStatusErr int
	FirstRequest  time.Time
	LastRequest   time.Time
	FirstResponse time.Time
	LastResponse  time.Time
}

var statOrder []string
var statData map[string]*Stat
var statChan chan Stat

func (s Stat) getTitles() []string {
	return []string{
		"Name", "Req Count", "Resp Count", "Resp Net Err", "Resp Status Err",
		"First Request", "Last Request", "First Response", "Last Response"}
}

func (s Stat) getStrings() []string {
	return []string{
		s.Name,
		fmt.Sprint(s.ReqCount),
		fmt.Sprint(s.RespCount),
		fmt.Sprint(s.RespNetErr),
		fmt.Sprint(s.RespStatusErr),
		fmt.Sprint(s.FirstRequest.Format("15:04:05.000")),
		fmt.Sprint(s.LastRequest.Format("15:04:05.000")),
		fmt.Sprint(s.FirstResponse.Format("15:04:05.000")),
		fmt.Sprint(s.LastResponse.Format("15:04:05.000")),
	}
}

func printStat(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(Stat{}.getTitles())

	for _, name := range statOrder {
		table.Append(statData[name].getStrings())
	}

	table.Render()
}

func statWriter() {
	var zeroTime time.Time
	for {
		s := <-statChan
		statData[s.Name].ReqCount += s.ReqCount
		statData[s.Name].RespCount += s.RespCount
		statData[s.Name].RespNetErr += s.RespNetErr
		statData[s.Name].RespStatusErr += s.RespStatusErr
		if s.LastRequest != zeroTime {
			statData[s.Name].LastRequest = s.LastRequest
			if statData[s.Name].FirstRequest == zeroTime {
				statData[s.Name].FirstRequest = s.LastRequest
			}
		}
		if s.LastResponse != zeroTime {
			statData[s.Name].LastResponse = s.LastResponse
			if statData[s.Name].FirstResponse == zeroTime {
				statData[s.Name].FirstResponse = s.LastResponse
			}
		}
	}
}

func resetStat(config Config) {
	for _, rule := range config.Rules {
		statData[rule.Request.Name] = &Stat{Name: rule.Request.Name}
	}
}

func initStat(config Config) {
	statOrder = make([]string, len(config.Rules))
	statData = make(map[string]*Stat)
	statChan = make(chan Stat, 1000)

	for i, rule := range config.Rules {
		statOrder[i] = rule.Request.Name
		statData[rule.Request.Name] = &Stat{Name: rule.Request.Name}
	}

	go statWriter()
}
