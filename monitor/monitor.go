package monitor

import (
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	monitor_controller "github.com/ProjectAthenaa/sonic-core/protos/monitorController"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/walmart/config"
	"github.com/prometheus/common/log"
)

var _ face.MonitorCallback = (*Task)(nil)

type Task struct {
	*base.BMonitor
	sku        string
	pxResponse []byte
	pxUUID string
}

func NewTask(data *monitor_controller.Task) (*Task, error) {
	task := &Task{
		BMonitor: &base.BMonitor{Data: data},
		sku:      data.Metadata[*config.Module.Fields[0].FieldKey],
	}
	task.Callback = task

	return task, nil
}

func (tk *Task) TaskLoop() {
	var err error
	for {
		select {
		case <-tk.Ctx.Done():
			return
		default:
			if err = tk.iteration(); err != nil {
				log.Error("error completing iteration: ", err)
			}
		}
	}
}

func (tk *Task) OnStarting() {
	tk.Client.Jar = fasttls.NewJar()
}

func (tk *Task) OnStopping() {
	//panic("")
}
