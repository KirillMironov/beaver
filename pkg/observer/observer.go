package observer

import (
	"fmt"
	"sync"
)

type Observer struct {
	messages []string
	mu       sync.RWMutex
}

func New() *Observer {
	return &Observer{}
}

func (o *Observer) Info(args ...any) {
	o.append(fmt.Sprint(args...))
}

func (o *Observer) Infof(format string, args ...any) {
	o.append(fmt.Sprintf(format, args...))
}

func (o *Observer) Error(args ...any) {
	o.append(fmt.Sprint(args...))
}

func (o *Observer) Errorf(format string, args ...any) {
	o.append(fmt.Sprintf(format, args...))
}

func (o *Observer) append(message string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.messages = append(o.messages, message)
}

func (o *Observer) First() string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if len(o.messages) == 0 {
		return ""
	}

	return o.messages[0]
}

func (o *Observer) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return len(o.messages)
}
