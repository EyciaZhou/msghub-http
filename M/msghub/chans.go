package msghub

import (
	"sync"
	"fmt"
	"time"
)

var (
	chans map[string]*ChanInfo	//TODO: move this part to redis?
	chansArray []*ChanInfo
	chansMutex sync.RWMutex
)

func updateChanInfos() (_err error) {
	defer func() {
		err := recover()
		if err != nil {
			_err = err.(error)
		}
	}()

	rows, err := db.Query(`
		SELECT
				id, Title, LastModify
			FROM topic`) //TODO: Maybe too much...

	if err != nil {
		return err
	}

	tmpChans := map[string]*ChanInfo{}
	tmpChansArray := []*ChanInfo{}

	var cnt int
	for cnt = 0; rows.Next(); cnt++ {
		topic := &ChanInfo{}

		_err = rows.Scan(&topic.Id, &topic.Title, &topic.LastModify)

		if _err != nil {
			return _err
		}

		tmpChans[topic.Id] = topic
		tmpChansArray = append(tmpChansArray, topic)
	}

	chansMutex.RLock()

	same := true
	if len(chans) != len(tmpChans) {
		same = false
	}
	for k, v := range chans {
		if vtmp, ok := tmpChans[k]; !ok {
			same = false
		} else {
			same = same && (v.Title==vtmp.Title) && (v.LastModify==vtmp.LastModify) && (v.Id==vtmp.Id)
		}
	} // because have same length, if they not same, chans can't include tmpchans

	chansMutex.RUnlock()

	if !same {
		chansMutex.Lock()
		chans = tmpChans
		chansArray = tmpChansArray
		chansMutex.Unlock()
	}

	return nil
}

func cronUpdateChanInfos() {
	time.Sleep(10*time.Second)

	for ;; {
		err := updateChanInfos()
		if err != nil {
			//TODO: LOG ERROR.... not decided which log module to use yet
			fmt.Printf("error: " + err.Error())
		}

		time.Sleep(30*time.Second)
	}
}

func init() {
	chans = map[string]*ChanInfo{}
	go cronUpdateChanInfos()
}
