package worker

import (
	"fmt"
	"gopkg.in/robfig/cron.v2"
)

// A way to keep track of user defined
// objects in cron pool.
type slave struct {
	id      int
	entryID cron.EntryID
	jobFunc func()
}

// Holds cron instance and tracks
// all jobs run in cron pool.
type Master struct {
	cronner *cron.Cron
	slaves  []*slave
}

// Defines a new Master and inits the cron package.
func Init() *Master {
	master := new(Master)
	master.cronner = cron.New()
	return master
}

// Start cron
func (m *Master) Start() {
	m.cronner.Start()
}

// Stop cron
func (m *Master) Stop() {
	m.cronner.Stop()
}

// Finds a slave that matches given id
func (m *Master) findSlave(id int) (*slave, error) {
	var foundslave *slave
	foundslave = nil
	for _, slave := range m.slaves {
		if slave.id == id {
			foundslave = slave
			break
		}
	}
	if foundslave == nil {
		return nil, fmt.Errorf("slave not found")
	}
	return foundslave, nil
}

// Removes slave from cron pool and Master's tracked slaves.
func (m *Master) RemoveSlave(id int) error {
	foundSlave, err := m.findSlave(id)
	if err != nil {
		return err
	}
	m.cronner.Remove(foundSlave.entryID)
	for index, slave := range m.slaves {
		if foundSlave.id == slave.id {
			m.slaves = append(m.slaves[:index], m.slaves[index+1:]...)
			break
		}
	}
	return nil
}

// Finds existing slave, removes it, and creates a new
// one with same id and new interval
func (m *Master) ChangeTime(id int, interval int) error {
	slave, err := m.findSlave(id)
	if err != nil {
		return err
	}
	err = m.RemoveSlave(id)
	if err != nil {
		return err
	}
	err = m.AddSlave(id, interval, slave.jobFunc)
	if err != nil {
		return err
	}
	return nil
}

// Run cron job associated with given slave id
func (m *Master) RunSlave(id int) error {
	foundSlave, err := m.findSlave(id)
	if err != nil {
		return err
	}
	entry := m.cronner.Entry(foundSlave.entryID)
	entry.Job.Run()
	return nil
}

// Creates a new slave and adds it to existing cron pool
// to be run on given interval (in minutes) and adds it to
// Master's slaves.
func (m *Master) AddSlave(id int, interval int, job func()) error {
	if slave, _ := m.findSlave(id); slave != nil {
		return fmt.Errorf("Slave with that id already exists.")
	}
	slave := new(slave)
	slave.id = id
	slave.jobFunc = job
	schedule := fmt.Sprintf("@every %vm", interval)
	entryId, err := m.cronner.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	slave.entryID = entryId
	m.slaves = append(m.slaves, slave)
	return nil
}
