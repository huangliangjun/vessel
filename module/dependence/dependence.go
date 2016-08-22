package dependence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containerops/vessel/models"
)

// ParsePipelineTemplate parse executor map from pipelineSpecTemplate
func ParsePipelineTemplate(template *models.PipelineSpecTemplate) (map[string]*models.Executor, error) {
	pipeline := template.MetaData
	stageSpec := template.Spec

	executorMap := make(map[string]*models.Executor, 0)
	executorListMap := make(map[string][]string, 0)
	pipeline.Stages = make([]string, 0, len(stageSpec))
	for _, stage := range stageSpec {
		if stage.Name == "" {
			return nil, errors.New("Stage has an empty name")
		}
		_, ok := executorMap[stage.Name]
		if ok {
			return nil, fmt.Errorf("Stage has repeat name: %v", stage.Name)
		}
		pipeline.Stages = append(pipeline.Stages, stage.Name)

		executor := &models.Executor{
			Info: stage,
			From: strings.Split(stage.Dependence, ","),
		}

		executorMap[stage.Name] = executor

		for _, from := range executor.From {
			executorList, ok := executorListMap[from]
			if !ok {
				executorList = make([]string, 0, 10)
			}
			executorList = append(executorList, stage.Name)
			executorListMap[from] = executorList
		}
	}
	return executorMap, checkDependenceValidity(executorListMap, executorMap)
}

func checkDependenceValidity(executorListMap map[string][]string, executorMap map[string]*models.Executor) error {
	if len(executorListMap[""]) == 0 {
		return errors.New("The first start stage list is empty")
	}

	//Check dependence stage name is exist
	for dependenceName := range executorListMap {
		if dependenceName == "" {
			continue
		}
		_, ok := executorMap[dependenceName]
		if !ok {
			return fmt.Errorf("Dependence stage name: %v is not exist", dependenceName)
		}
	}

	//Check dependence directed acyclic graph
	return checkEndlessChain(executorListMap, make([]string, 0, 10), "")

}

func checkEndlessChain(executorListMap map[string][]string, chain []string, checkName string) error {
	if checkName != "" {
		for _, chainItem := range chain {
			if chainItem == checkName {
				return fmt.Errorf("Dependence chain [%v,%v] is endless chain", strings.Join(chain, ","), checkName)
			}
		}
	}
	stageList, ok := executorListMap[checkName]
	if ok {
		for _, nextStage := range stageList {
			chain = append(chain, checkName)
			err := checkEndlessChain(executorListMap, chain, nextStage)
			if err != nil {
				return err
			}
			chain = chain[0 : len(chain)-1]
		}
	}
	return nil
}
