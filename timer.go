package tmx

import "time"

// timer can be used to calculate
// timestamp difference between GetElapsedTime and UpdateTime
type timer struct {
	beforeTime int64
	afterTime  int64
}

//Timer can be used to measure time delta between
//two snapshots
type Timer interface {
	GetElapsedTime() int64
	Start()
	UpdateTime()
}

//CreateTimer creates a new timer
func CreateTimer() Timer {
	return &timer{}
}

// Start the timer initially
func (t *timer) Start() {
	t.beforeTime = time.Now().UnixNano()
	t.afterTime = t.beforeTime
}

// GetElapsedTime returns the elapsed nano seconds
// since the last updateTime
func (t *timer) GetElapsedTime() int64 {
	elapsedTime := t.afterTime - t.beforeTime
	t.beforeTime = time.Now().UnixNano()

	return elapsedTime
}

//UpdateTime updates the afterTime to now
func (t *timer) UpdateTime() {
	t.afterTime = time.Now().UnixNano()
}
