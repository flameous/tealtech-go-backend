package tealtech

import (
	"sync"
)

type Database interface {
	GetUser(uid int, login string) *BotUser
	SaveUser(u *BotUser)
	Reset()
}

type DumpDatabase struct {
	mtx    *sync.RWMutex
	m      map[int]*BotUser
	logins map[string]int
}

func NewDumpDatabase() *DumpDatabase {
	d := DumpDatabase{
		mtx:    &sync.RWMutex{},
		m:      make(map[int]*BotUser),
		logins: make(map[string]int),
	}

	return &d
}

func (d *DumpDatabase) GetUser(uid int, login string) *BotUser {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	if login != "" {
		id, ok := d.logins[login]
		if !ok {
			return nil
		}
		uid = id
	}
	u, ok := d.m[uid]
	if !ok {
		return nil
	}
	return u
}

func (d *DumpDatabase) SaveUser(u *BotUser) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if u.Login != "" {
		d.logins[u.Login] = u.UserId
	}
	d.m[u.UserId] = u
}

func (d *DumpDatabase) Reset() {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	d.m = make(map[int]*BotUser)
	d.logins = make(map[string]int)
}
