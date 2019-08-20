package kwidgets

import (
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"

	"ktop/app"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	runningStyle   = ui.NewStyle(ui.ColorGreen)
	completedStyle = ui.NewStyle(ui.ColorBlue)
)

type PodList struct {
	*widgets.Table
	header         []string
	updateInterval time.Duration
	pods           []app.PodDescriptor
}

func NewPodList() *PodList {
	pl := &PodList{
		Table:          widgets.NewTable(),
		updateInterval: time.Second,
	}
	pl.Title = " Pods "
	pl.header = []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}
	pl.RowSeparator = false
	pl.ColumnResizer = func() {
		nameWidth, length := 0, 0
		for _, p := range pl.pods {
			length = len(p.Name)
			if length > nameWidth {
				nameWidth = length
			}
		}
		pl.ColumnWidths = []int{100, 8, 12, 9, 7}
	}

	pl.update()
	go func() {
		for range time.NewTicker(pl.updateInterval).C {
			pl.Lock()
			pl.update()
			pl.Unlock()
		}
	}()

	return pl
}

func (pl *PodList) update() {
	pl.pods = app.GetPods("default")
	pl.convertProcsToTableRows()
}

func (pl *PodList) convertProcsToTableRows() {
	strings := make([][]string, len(pl.pods))
	for i := range pl.pods {
		strings[i] = []string{
			pl.pods[i].Name,
			pl.pods[i].GetPodReadiness(),
			fmt.Sprintf(pl.getRowStyle(pl.pods[i].Status), pl.pods[i].Status),
			strconv.Itoa(int(pl.pods[i].Restart)),
			"age",
		}
	}
	pl.Rows = append([][]string{pl.header}, strings...)
}

func (pl *PodList) getRowStyle(status string) string {
	switch status {
	case string(v1.PodRunning):
		return "[%s](fg:blue)"
	case string(v1.PodSucceeded):
		return "[%s](fg:green)"
	case string(v1.PodFailed):
		return "[%s](fg:red)"
	case string(v1.PodPending):
		return "[%s](fg:orange)"
	default:
		return "%s"
	}
}
