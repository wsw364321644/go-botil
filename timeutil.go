package botil

import (
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

const ServerTimeFormat="2006-01-02T15:04:05.999Z"
const BriefTimeFormat="2006-01-02"

var CronMinuteParser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow ,
)



// Crontab crontab manager
type Crontab struct {
	inner *cron.Cron
	ids   map[string]cron.EntryID
	mutex *sync.RWMutex
}

// NewCrontab new crontab
func NewCrontab(mutex *sync.RWMutex ,opts ...cron.Option) *Crontab {
	return &Crontab{
		inner: cron.New(opts...),
		ids:   make(map[string]cron.EntryID),
		mutex:mutex,
	}
}

// IDs ...
func (c *Crontab) IDs() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	validIDs := make([]string, 0, len(c.ids))
	invalidIDs := make([]string, 0)
	for sid, eid := range c.ids {
		if e := c.inner.Entry(eid); e.ID != eid {
			invalidIDs = append(invalidIDs, sid)
			continue
		}
		validIDs = append(validIDs, sid)
	}
	for _, id := range invalidIDs {
		delete(c.ids, id)
	}
	return validIDs
}

// Start start the crontab engine
func (c *Crontab) Start() {
	c.inner.Start()
}

// Stop stop the crontab engine
func (c *Crontab) Stop() {
	c.inner.Stop()
}

// DelByID remove one crontab task
func (c *Crontab) DelByID(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	eid, ok := c.ids[id]
	if !ok {
		return
	}
	c.inner.Remove(eid)
	c.inner.Location()
	delete(c.ids, id)
}

// AddByID add one crontab task
// id is unique
// spec is the crontab expression
func (c *Crontab) AddJob(id string, spec string, cmd cron.Job) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if eid, ok := c.ids[id]; ok {
		c.inner.Remove(eid)
	}
	eid, err := c.inner.AddJob(spec, cmd)
	if err != nil {
		return err
	}
	c.ids[id] = eid
	return nil
}

// AddByFunc add function as crontab task
func (c *Crontab) AddFunc(id string, spec string, f func()) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if eid, ok := c.ids[id]; ok {
		c.inner.Remove(eid)
	}
	eid, err := c.inner.AddFunc(spec, f)
	if err != nil {
		return err
	}
	c.ids[id] = eid
	return nil
}

type defaultJob struct{
	f func()
}
func (self defaultJob)Run(){
	self.f()
}

func (c *Crontab) ScheduleFunc(id string, schedule cron.Schedule, f func()) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if eid, ok := c.ids[id]; ok {
		c.inner.Remove(eid)
	}

	eid:= c.inner.Schedule(schedule,defaultJob{f})
	c.ids[id] = eid
	return nil
}
func (c *Crontab) Schedule(id string, schedule cron.Schedule,job cron.Job) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if eid, ok := c.ids[id]; ok {
		c.inner.Remove(eid)
	}

	eid:= c.inner.Schedule(schedule,job)
	c.ids[id] = eid
	return nil
}

// IsExists check the crontab task whether existed with job id
func (c *Crontab) IsExists(jid string) bool {
	_, exist := c.ids[jid]
	return exist
}

func (c *Crontab) IsValid(jid string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	id := c.ids[jid]
	return c.inner.Entry(id).Valid()

}
func (c *Crontab) GetNext(jid string) time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	id := c.ids[jid]
	return c.inner.Entry(id).Next
}