package fetch

import (
	"fmt"
	"sync"

	"github.com/ultimatums/ultimate/config"
	"github.com/upmio/horus/log"
)

type Task interface {
	Identity() []string

	Run(ch chan<- *config.UnitConfig)

	Stop()
}

type TaskManager struct {
	m       sync.RWMutex
	running bool

	// Units by their ID.
	units map[string]*Unit
	// Tasks by the task config they are derived from.
	tasks map[*config.TaskConfig]Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		units: make(map[string]*Unit),
	}
}

func (tm *TaskManager) Run() {
	log.Info("Starting task manager...")

	identities := map[string]struct{}{}

	for taskCfg, task := range tm.tasks {
		log.Infof("taskCfg = %v", taskCfg)
		ch := make(chan *config.UnitConfig)
		// Every task has a batch of units.
		go tm.handleUnitUpdates(taskCfg, ch)

		for _, id := range task.Identity() {
			id = fullId(taskCfg, id)
			identities[id] = struct{}{}
		}

		defer func(t Task, c chan *config.UnitConfig) {
			go t.Run(c)
		}(task, ch)
	}

	log.Info("before removeUnits")

	tm.m.Lock()
	defer tm.m.Unlock()

	tm.removeUnits(func(id string) bool {
		if _, ok := identities[id]; ok {
			return false
		}
		return true
	})

	tm.running = true
	log.Info("Run finished.")
}

func fullId(cfg *config.TaskConfig, identity string) string {
	return cfg.TaskName + ":" + identity
}

func (tm *TaskManager) handleUnitUpdates(taskCfg *config.TaskConfig, ch <-chan *config.UnitConfig) {
	for unitCfg := range ch {
		//		log.Debugf("Received potential update for unit config %s", unitCfg.Identity)
		err := tm.updateUnitSet(unitCfg, taskCfg)
		if err != nil {
			log.Errorf("Error updating units: %s", err)
		}
	}
	log.Debug("handleUnitUpdates finished.")
}

func (tm *TaskManager) updateUnitSet(unitCfg *config.UnitConfig, taskCfg *config.TaskConfig) error {

	id := fullId(taskCfg, unitCfg.Identity)
	log.Info("id = ", id)

	newUnit := NewUnit(unitCfg)

	tm.m.Lock()
	defer tm.m.Unlock()

	if !tm.running {
		return nil
	}

	oldUnit, ok := tm.units[id]
	if ok {
		isMatch := (oldUnit.identity == newUnit.identity)
		log.Debug("is match = ", isMatch)
		if isMatch {
			//TODO oldUnit.Update()
		} else {
			//TODO newUnit.Start()
			//TODO oldUnit.Stop()
			delete(tm.units, id)
		}
	} else {
		//TODO newUnit.Start()
		tm.units[id] = newUnit
	}
	log.Info("units = ", tm.units)

	return nil
}

func (tm *TaskManager) Stop() {
	tm.m.RLock()
	if tm.running {
		defer tm.stop(true)
	}
	defer tm.m.RUnlock()
}

func (tm *TaskManager) stop(removeUnits bool) {
	log.Info("Stoping task manager...")
	defer log.Info("Task manager stopped.")

	tm.m.Lock()
	tasks := []Task{}
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	tm.m.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		go func(ts Task) {
			ts.Stop()
			wg.Done()
		}(task)
	}
	wg.Wait()

	tm.m.Lock()
	defer tm.m.Unlock()

	if removeUnits {
		tm.removeUnits(nil)
	}

	tm.running = false
}

func (tm *TaskManager) removeUnits(f func(string) bool) {
	if f == nil {
		f = func(string) bool { return true }
	}

	log.Debug("length of tm.units is ", len(tm.units))

	var wg sync.WaitGroup
	for id, unit := range tm.units {
		if !f(id) {
			continue
		}
		wg.Add(1)
		go func(u *Unit) {
			//TODO
			wg.Done()
		}(unit)
	}
	wg.Wait()
}

func (tm *TaskManager) ApplyConfig(cfg *config.Config) {
	tm.m.RLock()
	running := tm.running
	tm.m.RUnlock()

	if running {
		tm.Stop()
		defer tm.Run()
	}

	tasks := map[*config.TaskConfig]Task{}
	for _, taskCfg := range cfg.TaskConfigs {
		if taskCfg.FetchInterval == 0 {
			taskCfg.FetchInterval = cfg.GlobalConfig.FetchInterval
		}
		tasks[taskCfg] = BuildTaskFromConfig(taskCfg)
	}

	tm.m.Lock()
	defer tm.m.Unlock()
	tm.tasks = tasks
}

func BuildTaskFromConfig(taskCfg *config.TaskConfig) Task {
	for i, unitCfg := range taskCfg.UnitConfigs {
		if unitCfg.FetchInterval == 0 {
			unitCfg.FetchInterval = taskCfg.FetchInterval
		}
		unitCfg.Identity = fmt.Sprintf("sample:%d", i)
	}

	var task Task
	task = NewSampleTask(taskCfg.UnitConfigs)
	return task
}

type SampleTask struct {
	UnitConfigs []*config.UnitConfig
}

func NewSampleTask(unitConfigs []*config.UnitConfig) *SampleTask {
	return &SampleTask{
		UnitConfigs: unitConfigs,
	}
}

func (t *SampleTask) Run(ch chan<- *config.UnitConfig) {
	log.Info("length of unitConfigs = ", len(t.UnitConfigs))
	for _, unitCfg := range t.UnitConfigs {
		log.Info("unitCfg: ", unitCfg)
		ch <- unitCfg
	}
	close(ch)
}

func (t *SampleTask) Identity() []string {
	ret := make([]string, len(t.UnitConfigs))
	for i, unitCfg := range t.UnitConfigs {
		ret[i] = unitCfg.Identity
	}
	return ret
}

func (t *SampleTask) Stop() {
}
