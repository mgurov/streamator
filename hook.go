package main

import (
	"sync"
	"github.com/sirupsen/logrus"
	"fmt"
)

type cappedInMemoryRecorderHook struct {
	m sync.Mutex
	records []*logrus.Entry
	wIndex  int
	owerwrites bool
}

func newCappedInMemoryRecorderHook(c int) *cappedInMemoryRecorderHook {
	if c <= 0 {
		panic("cappedInMemoryRecorderHook should be not empty size of")
	}
	return &cappedInMemoryRecorderHook{
		records: make([]*logrus.Entry, c),
	}
}

func (h *cappedInMemoryRecorderHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *cappedInMemoryRecorderHook) Fire(e *logrus.Entry) error {
	h.m.Lock()
	defer h.m.Unlock()

	h.records[h.wIndex] = e

	h.wIndex ++
	if h.wIndex >= len(h.records) {
		h.wIndex = 0
		h.owerwrites = true
	}
	return nil
}

func (h *cappedInMemoryRecorderHook) Copy() []*logrus.Entry {
	h.m.Lock()
	defer h.m.Unlock()

	if (!h.owerwrites) {
		//todo: test me now
		return h.records[:h.wIndex]
	}

	fmt.Println("Copying", h.records)
	fmt.Println("wIndex", h.wIndex)

	result := make([]*logrus.Entry, len(h.records))

	fmt.Println("h.records[h.wIndex:]", h.records[h.wIndex:])
	copy(result, h.records[h.wIndex:])
	fmt.Println("result", result)

	fmt.Println("result[h.wIndex:]", result[len(result) - h.wIndex:])
	fmt.Println("h.records[:h.wIndex]", h.records[:h.wIndex])

	copy(result[len(result) - h.wIndex:], h.records[:h.wIndex])
	fmt.Println("result", result)

	return result
}
