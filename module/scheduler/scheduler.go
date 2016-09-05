package scheduler

import (
	"fmt"
	"log"

	"sync"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/stage"
)

var wg sync.WaitGroup

type schedulerHand func(interface{}, map[string]bool, chan *models.ExecutedResult)

// StartStage start point on scheduler
func StartPoint(executorList []interface{}, startMark string) []*models.ExecutedResult {
	return execute(executorList, startMark, stage.Start)
}

// StopPoint stop point on scheduler
func StopPoint(executorList []interface{}, startMark string) []*models.ExecutedResult {
	return StopExecute(executorList, startMark, stage.Stop)
}

func execute(executorList []interface{}, startMark string, handler schedulerHand) []*models.ExecutedResult {
	count := len(executorList)
	readyMap := map[string]bool{startMark: true}
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)
	running := true
	for running {
		for _, executor := range executorList {
			wg.Add(1)
			go func(exec interface{}, ready map[string]bool, finish chan *models.ExecutedResult) {
				defer wg.Done()
				handler(exec, ready, finish)
			}(executor, readyMap, finishChan)
		}
		result := <-finishChan

		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			sv := &models.StageVersion{
				SID:    result.SID,
				Detail: result.Detail,
			}
			err := sv.Update()
			if err != nil {
				return resultList
			}
			running = false
		} else {
			resultLen := len(resultList)
			running = resultLen != count
			log.Println(fmt.Sprintf("scheduler StartStage name = %v and finish num = %d", result.Name, resultLen))
		}
	}

	wg.Wait()
	return resultList
}

func StopExecute(executorList []interface{}, startMark string, handler schedulerHand) []*models.ExecutedResult {
	count := len(executorList)
	readyMap := map[string]bool{startMark: true}
	finishChan := make(chan *models.ExecutedResult, count)
	resultList := make([]*models.ExecutedResult, 0, count)
	for _, executor := range executorList {
		wg.Add(1)
		go func(exec interface{}, ready map[string]bool, finish chan *models.ExecutedResult) {
			defer wg.Done()
			handler(exec, ready, finish)
		}(executor, readyMap, finishChan)
	}

	for i := 0; i < count; i++ {
		result := <-finishChan
		resultList = append(resultList, result)
		if result.Status != models.ResultSuccess {
			sv := &models.StageVersion{
				SID:    result.SID,
				Detail: result.Detail,
			}
			err := sv.Update()
			if err != nil {
				return resultList
			}
		} else {
			resultLen := len(resultList)
			log.Println(fmt.Sprintf("scheduler StartStage name = %v and finish num = %d", result.Name, resultLen))
		}
	}

	wg.Wait()
	return resultList
}
