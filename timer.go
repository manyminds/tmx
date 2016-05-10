package tmx

import "time"

// Timer can be used to calculate
// timestamp difference between GetElapsedTime and UpdateTime
type timer struct {
	beforeTime int64
	afterTime  int64
}

func createTimer() *timer {
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
