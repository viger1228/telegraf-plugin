// File: tcping.go
// Author: walker
// Mail: walkerIVI@gmail.com
// Changelogs:
//   2018.10.26: init

package tcping

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/viger1228/golib/tcping"
)

type Tcping struct {
}

func (self *Tcping) SampleConfig() string {
	return ""
}

func (self *Tcping) Description() string {
	return ""
}

func (self *Tcping) Gather(acc telegraf.Accumulator) error {

	target := []string{
		"www.baidu.com:80",
		"www.google.com:80",
		"web.telegram.org:443",
	}

	var wg sync.WaitGroup

	for _, t := range target {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			tags := make(map[string]string)
			fields := make(map[string]interface{})

			url := strings.Split(t, ":")
			port, _ := strconv.Atoi(url[1])
			cli := tcping.TCPinger{
				Target:   url[0],
				Port:     port,
				Times:    10,
				Timeout:  2,
				Interval: 1,
				Statis:   map[string]float64{},
			}
			cli.Run()

			tags["target"] = cli.Target
			tags["ip"] = cli.IP
			tags["port"] = fmt.Sprintf("%v", cli.Port)

			fields["num"] = cli.Statis["num"]
			fields["max"] = cli.Statis["max"]
			fields["min"] = cli.Statis["min"]
			fields["avg"] = cli.Statis["avg"]
			fields["std"] = cli.Statis["std"]
			fields["loss"] = cli.Statis["loss"]

			acc.AddFields("tcping", fields, tags)

		}(t)
	}
	wg.Wait()
	return nil
}

func init() {
	inputs.Add("tcping", func() telegraf.Input { return &Tcping{} })
}
