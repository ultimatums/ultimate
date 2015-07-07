package fetch

import (
	"fmt"
	"sync"

	"github.com/ultimatums/ultimate/config"
	"github.com/upmio/horus/log"
)

type Task interface {
	Identity() []string

	Run(ch chan<- *config.UnitSet)

	Stop()
}

type TaskManager struct {
	m       sync.RWMutex
	running bool

	// Units by their ID.
	units map[string][]*Unit
	// Tasks by the task config they are derived from.
	tasks map[*config.TaskConfig][]Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		units: make(map[string][]*Unit),
	}
}

func (tm *TaskManager) Run() {
	log.Info("Starting task manager...")

	for taskCfg, tasks := range tm.tasks {
		for _, task := range tasks {
			ch := make(chan *config.UnitSet)
			go tm.handleUnitUpdates(taskCfg, ch)

			defer func(t Task, c chan *config.UnitSet) {
				go task.Run(c)
			}(task, ch)
		}
	}

	tm.m.Lock()
	defer tm.m.Unlock()

	tm.running = true
}

func (tm *TaskManager) Stop() {
	log.Info("Stoping task manager...")
}

func (tm *TaskManager) handleUnitUpdates(taskCfg *config.TaskConfig, ch <-chan *config.UnitSet) {
	for unitSet := range ch {
		log.Debugf("Received potential update for unit set %s", unitSet.Identity)
		err := tm.updateUnitSet(unitSet, taskCfg)
		if err != nil {
			log.Errorf("Error updating units: %s", err)
		}
	}
}

func (tm *TaskManager) updateUnitSet(unitSet *config.UnitSet, taskCfg *config.TaskConfig) error {
	tm.unitsFromSet(unitSet, taskCfg)
	return nil
}

func (tm *TaskManager) unitsFromSet(unitSet *config.UnitSet, taskCfg *config.TaskConfig) ([]*Unit, error) {
	tm.m.RLock()
	defer tm.m.RUnlock()

	for tagName, tagValue := range unitSet.UnitTags {
		log.Debugln("tagName = ", tagName, ", tagValue = ", tagValue)
	}

	return nil, nil
}

func (tm *TaskManager) ApplyConfig(cfg *config.Config) {
	tm.m.RLock()
	running := tm.running
	tm.m.RUnlock()

	if running {
		tm.Stop()
		defer tm.Run()
	}

	tasks := map[*config.TaskConfig][]Task{}
	for _, taskCfg := range cfg.FetchConfigs {
		tasks[taskCfg] = BuildTaskFromConfig(taskCfg)
	}

	tm.m.Lock()
	defer tm.m.Unlock()
	tm.tasks = tasks
}

func BuildTaskFromConfig(taskCfg *config.TaskConfig) []Task {
	var tasks []Task
	switch taskCfg.TaskName {
	case "host", "container":
		tasks = append(tasks, NewSampleTask(taskCfg.UnitSets))
	}
	return tasks
}

type SampleTask struct {
	UnitSets []*config.UnitSet
}

func NewSampleTask(unitSets []*config.UnitSet) *SampleTask {
	for i, unitSet := range unitSets {
		unitSet.Identity = fmt.Sprintf("sample:%d", i)
	}
	return &SampleTask{
		UnitSets: unitSets,
	}
}

func (t *SampleTask) Run(ch chan<- *config.UnitSet) {
	for _, unitSet := range t.UnitSets {
		ch <- unitSet
	}
	close(ch)
}

func (t *SampleTask) Identity() []string {
	ret := make([]string, len(t.UnitSets))
	for i, unitSet := range t.UnitSets {
		ret[i] = unitSet.Identity
	}
	return ret
}

func (t *SampleTask) Stop() {
}
