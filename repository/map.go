package repository

import "sync"

type MapWait struct {
	wait    map[string]int
	sync.Mutex
}

type MapResp struct {
	resp    map[string][]byte
	sync.Mutex

}

func (wait *MapWait) addWait(req string) {
	wait.Lock()
	defer wait.Unlock()

	wait.wait[req]++
}

func (wait *MapWait) deleteWait(req string) {
	wait.Lock()
	defer wait.Unlock()
	delete(wait.wait, req)
}

func (wait *MapWait) checkSing(req string) (int, bool) {
	wait.Lock()
	defer wait.Unlock()
	val, ok := wait.wait[req]
	return val, ok
}


func (wait *MapResp) addWait(req string, data []byte) {
	wait.Lock()
	defer wait.Unlock()

	wait.resp[req] = data
}

func (wait *MapResp) deleteWait(req string) {
	wait.Lock()
	defer wait.Unlock()
	delete(wait.resp, req)
}

func (wait *MapResp) checkSing(req string) ([]byte, bool) {
	wait.Lock()
	defer wait.Unlock()
	val, ok := wait.resp[req]
	return val, ok
}
