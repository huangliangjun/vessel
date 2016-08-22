package kubernetes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils/timer"
)

// GetBusinessRes get business result with kubernetes pods
func GetBusinessRes(stage *models.Stage, hourglass *timer.Hourglass) (res *models.K8SRes) {
	if hourglass.GetLeftNanoseconds() < 0 {
		return formatResult(models.ResultTimeout, "Get business result in kubernetes timeout")
	}
	checkCount := stage.StatusCheckCount
	checkInterval := stage.StatusCheckInterval
	if checkInterval == 0 {
		checkInterval = 30
	}
	if checkCount == 0 {
		checkCount = 3
	}

	resPods := make(map[string]int)
	ipList, err := getPodIPList(stage)
	if err != nil {
		return formatResult(models.ResultFailed, err.Error())
	}

	checkCh := make(chan bool)
	for _, item := range ipList {
		checkURL := fmt.Sprintf("http://%v:%v%v", item, stage.Port, stage.StatusCheckURL)
		go getPodResult(checkURL, checkCount, checkInterval, resPods, checkCh, hourglass)
	}
	for i := 0; i < len(ipList); i++ {
		<-checkCh
	}
	close(checkCh)
	return formatBusResult(resPods)
}

func getPodResult(checkURL string, count uint64, interval uint64, resPods map[string]int, checkCh chan bool, hourglass *timer.Hourglass) {
	resCode := -1
	hasTime := hourglass.GetLeftNanoseconds() > 0
	for hasTime {
		resCode = httpGet(checkURL)
		if resCode == 200 || resCode == 0 {
			resPods[checkURL] = 200
			checkCh <- true
			return
		}
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			if count--; count == 0 {
				hasTime = false
			}
		case <-time.After(hourglass.GetLeftNanoseconds()):
			hasTime = false
		}
	}
	resPods[checkURL] = resCode
	checkCh <- false
}

func httpGet(checkURL string) int {
	resp, err := http.Get(checkURL)
	if err != nil {
		return -1
	}
	_, err = ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		resp.Body.Close()
	}
	if err != nil {
		return -1
	}
	return resp.StatusCode
}

func formatBusResult(mapRes map[string]int) *models.K8SRes {
	detail := ""
	for key, value := range mapRes {
		if value == 200 {
			return formatResult(models.ResultSuccess, "")
		}
		if detail == "" {
			detail = "Get business from " + key
		} else {
			detail += " and " + key
		}
	}
	return formatResult(models.ResultTimeout, detail)
}
