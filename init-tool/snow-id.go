package init_tool

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func SnowIDInit() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", Conf.ProjectConfig.StartTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(Conf.ProjectConfig.MachineID)
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}
