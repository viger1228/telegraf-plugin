// File: tcping.go
// Author: walker
// Mail: walkerIVI@gmail.com
// Changelogs:
//   2018.10.26: init

package tcping

import (
	"fmt"
	"os"
	//"strconv"
	//"strings"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/viger1228/golib/mysql"
	"github.com/viger1228/golib/tcping"
	"github.com/viger1228/golib/tool"
)

type Tcping struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Times    int
	Timeout  int
}

func (self *Tcping) SampleConfig() string {
	return ""
}

func (self *Tcping) Description() string {
	return ""
}

func (self *Tcping) Gather(acc telegraf.Accumulator) error {

	//target := []string{
	//	"www.baidu.com:80",
	//	"www.google.com:80",
	//	"web.telegram.org:443",
	//}

	hostname, err := os.Hostname()
	tool.CheckErr(err)
	sql := fmt.Sprintf("SELECT `hostname`,`target`,`port` FROM "+
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
				Interval: 1,
				Statis:   map[string]float64{},
			}
			cli.Run()

			tags["@target"] = cli.Target
			tags["@ip"] = cli.IP
			tags["@port"] = fmt.Sprintf("%v", cli.Port)

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
