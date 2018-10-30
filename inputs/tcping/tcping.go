// File: tcping.go
// Author: walker
// Mail: walkerIVI@gmail.com
// Changelogs:
//   2018.10.26: init

package tcping

import (
	"fmt"
	"os"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/viger1228/golib/mysql"
	"github.com/viger1228/golib/tcping"
)

const sampleConfig = `
  ## Task Database
  host = "127.0.0.1"
  port = 3306
  user = ""
  password = ""
  database = "mon"
  ## Execute Time
  times = 20
  ## Execute Timeout
  timeout = 2
  ## Execute Interval
  interval = 2
`

type Tcping struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Times    int
	Timeout  int
	Interval int
}

func (self *Tcping) SampleConfig() string {
	return sampleConfig
}

func (self *Tcping) Description() string {
	return "Tcping from DB Task list and return statistics"
}

func (self *Tcping) Gather(acc telegraf.Accumulator) error {

	//target := []string{
	//	"www.baidu.com:80",
	//	"www.google.com:80",
	//	"web.telegram.org:443",
	//}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("SELECT `hostname`,`type`,`target`,`port`,`note` FROM "+
		"t_telegraf_tcping WHERE `enable`=1 AND `hostname`='%v'", hostname)

	api := mysql.MySQL{
		Host:     self.Host,
		Port:     self.Port,
		User:     self.User,
		Password: self.Password,
		Database: self.Database,
	}

	target := api.Query(sql)

	var wg sync.WaitGroup

	for _, t := range target {
		wg.Add(1)
		go func(t map[string]interface{}) {
			defer wg.Done()
			tags := make(map[string]string)
			fields := make(map[string]interface{})

			cli := tcping.TCPinger{
				Target:   t["target"].(string),
				Port:     t["port"].(int),
				Times:    self.Times,
				Timeout:  self.Timeout,
				Interval: self.Interval,
				Statis:   map[string]float64{},
			}
			cli.Run()

			tags["@target"] = cli.Target
			tags["@type"] = t["type"].(string)
			tags["@ip"] = cli.IP
			tags["@port"] = fmt.Sprintf("%v", cli.Port)
			tags["@note"] = t["note"].(string)

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
