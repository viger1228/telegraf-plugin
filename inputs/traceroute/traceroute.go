// File: traceroute.go
// Author: walker
// Mail: walkerIVI@gmail.com
// Changelogs:
//   2018.10.26: init

package traceroute

import (
	"fmt"
	"os"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/viger1228/golib/mysql"
	"github.com/viger1228/golib/traceroute"
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
  timeout = 2
`

type Traceroute struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Times    int
	Timeout  int
	Interval int
}

func (self *Traceroute) SampleConfig() string {
	return sampleConfig
}

func (self *Traceroute) Description() string {
	return "Traceroute from DB Task list and return statistics"
}

func (self *Traceroute) Gather(acc telegraf.Accumulator) error {

	//target := []map[string]interface{}{}
	//data := map[string]interface{}{
	//	"target": "47.91.20.102",
	//}
	//target = append(target, data)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("SELECT `hostname`,`type`,`target`,`note` FROM "+
		"t_telegraf_traceroute WHERE `enable`=1 AND `hostname`='%v'", hostname)

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

			cli := traceroute.Tracer{
				Target:   t["target"].(string),
				Times:    self.Times,
				Timeout:  self.Timeout,
				Interval: self.Interval,
			}
			cli.Run()

			for _, host := range cli.Statis {
				tags := make(map[string]string)
				fields := make(map[string]interface{})

				tags["@target"] = cli.Target
				tags["@type"] = t["type"].(string)
				tags["@hop"] = fmt.Sprintf("%02d_%v", host.Hop, host.IP)
				tags["@note"] = t["note"].(string)

				fields["num"] = host.Num
				fields["max"] = host.Max
				fields["min"] = host.Min
				fields["avg"] = host.Avg
				fields["std"] = host.Std
				fields["loss"] = host.Loss

				acc.AddFields("traceroute", fields, tags)
			}

		}(t)

	}
	wg.Wait()
	return nil
}

func init() {
	inputs.Add("traceroute", func() telegraf.Input { return &Traceroute{} })
}
