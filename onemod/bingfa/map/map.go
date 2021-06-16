/*
 *@Description
 *@author          lirui
 *@create          2021-06-15 19:27
 */
package main

import "sync"

type RWMap struct {
	sync.RWMutex
	m map[int]int
}

func NewRWMap(n int) *RWMap {
	return &RWMap{
		m: make(map[int]int, n),
	}
}

// 对map的增删查改进行在操作前和操作后进行加锁
func (m *RWMap) Get(k int) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k]
	return v, existed
}

func main() {

}
