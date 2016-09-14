package db_test

import (
	//"errors"
	"testing"

	"github.com/containerops/vessel/db"
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
)

func Test_InitDB(t *testing.T) {
	if err := setting.InitGlobalConf("../conf/global.yaml"); err != nil {
		t.Fatal("Read config failure : ", err.Error())
	}
	if err := db.InitDB(setting.RunTime.Database.Driver, setting.RunTime.Database.Username,
		setting.RunTime.Database.Password, setting.RunTime.Database.Host+":"+setting.RunTime.Database.Port,
		setting.RunTime.Database.Schema); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(models.Pipeline), new(models.PipelineVersion)); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(models.Stage), new(models.StageVersion)); err != nil {
		t.Fatal(err)
	}
	if err := db.Instance.RegisterModel(new(models.Point), new(models.PointVersion)); err != nil {
		t.Fatal(err)
	}
}

func Test_Create(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if err := db.Instance.Create(pipeline); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	pipelineVsn := &models.PipelineVersion{
		PID:   pipeline.ID,
		State: "running",
	}
	if err := db.Instance.Create(pipelineVsn); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	stage := &models.Stage{
		PID:           pipeline.ID,
		Namespace:     "test",
		Name:          "master",
		Replicas:      1,
		Type:          "container",
		Dependencies:  "",
		ArtifactsJSON: "aaa",
		VolumesJSON:   "aaa",
		PortsJSON:     "aaa",
	}
	if err := db.Instance.Create(stage); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	stage = &models.Stage{
		PID:           pipeline.ID,
		Namespace:     "test",
		Name:          "slave",
		Replicas:      1,
		Type:          "container",
		Dependencies:  "",
		ArtifactsJSON: "bbb",
		VolumesJSON:   "bbb",
		PortsJSON:     "bbb",
	}
	if err := db.Instance.Create(stage); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	stageVsn := &models.StageVersion{
		PvID:  pipelineVsn.ID,
		SID:   stage.ID,
		State: "running",
	}
	if err := db.Instance.Create(stageVsn); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	point := &models.Point{
		PID:        pipeline.ID,
		Type:       "start",
		Triggers:   "a,b,c",
		Conditions: "a,b,c",
	}
	if err := db.Instance.Create(point); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	point = &models.Point{
		PID:        pipeline.ID,
		Type:       "end",
		Triggers:   "a,b,c",
		Conditions: "a,b,c",
	}
	if err := db.Instance.Create(point); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	pointVsn := &models.PointVersion{
		PvID:    pipelineVsn.ID,
		PointID: point.ID,
		State:   "running",
	}
	if err := db.Instance.Create(pointVsn); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	db.Instance.Commit()
}

func Test_Count(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		t.Fatal(err)
	} else if pipeline.ID <= 0 {
		t.Fatal("query error")
	}
	t.Log(pipeline)
	stage := &models.Stage{
		PID:       pipeline.ID,
		Namespace: "test",
		Name:      "master",
	}
	if _, err := db.Instance.Count(stage); err != nil {
		t.Fatal(err)
	}
	t.Log(stage)

}

func Test_Save(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		t.Fatal(err)
	} else if pipeline.ID <= 0 {
		t.Fatal("query error")
	}
	pipeline2 := &models.Pipeline{
		ID:        pipeline.ID,
		Namespace: "test",
		Name:      "test",
		Timeout:   60,
		CreatedAt: pipeline.CreatedAt,
	}
	if err := db.Instance.Save(pipeline2); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	t.Log(pipeline2)
	db.Instance.Commit()
}

func Test_Update(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		t.Fatal(err)
	} else if pipeline.ID <= 0 {
		t.Fatal("query error")
	}
	pipeline2 := &models.Pipeline{
		ID:        pipeline.ID,
		Namespace: "test",
		Name:      "test",
		Timeout:   120,
	}
	if err := db.Instance.Update(pipeline2); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	t.Log(pipeline2)
	db.Instance.Commit()
}

func Test_UpdateField(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		t.Fatal(err)
	} else if pipeline.ID <= 0 {
		t.Fatal("query error")
	}
	pipeline2 := &models.Pipeline{
		ID:        pipeline.ID,
		Namespace: "test",
		Name:      "test",
		Timeout:   120,
	}
	if err := db.Instance.UpdateField(pipeline2, "Timeout", 180); err != nil {
		t.Fatal(err)
	}
	t.Log(pipeline2)
}

func Test_QueryM(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test1",
		Name:      "test1",
		Timeout:   30,
	}
	if err := db.Instance.Create(pipeline); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	db.Instance.Commit()
	//get namespace == test1
	results := []models.Pipeline{}
	if err := db.Instance.QueryM(&models.Pipeline{Namespace: "test1"}, &results); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}

	//get all
	results1 := []models.Pipeline{}
	if err := db.Instance.QueryM(&models.Pipeline{}, &results1); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results1)
	}
}

func Test_QueryF(t *testing.T) {
	//get namespace like te
	results := []models.Pipeline{}
	if err := db.Instance.QueryF(&models.Pipeline{Namespace: "te"}, &results); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}
	//get namespace like te AND name like 1
	results1 := []models.Pipeline{}
	if err := db.Instance.QueryF(&models.Pipeline{Namespace: "te", Name: "1"}, &results1); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results1)
	}
}

func Test_SoftDelete(t *testing.T) {
	pipeline := &models.Pipeline{
		Namespace: "test",
		Name:      "test",
	}
	if _, err := db.Instance.Count(pipeline); err != nil {
		t.Fatal(err)
	} else if pipeline.ID <= 0 {
		t.Fatal("query error")
	}
	if err := db.Instance.DeleteS(&models.Pipeline{ID: pipeline.ID}); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	db.Instance.Commit()
	if err := db.Instance.DeleteS(&models.Stage{PID: pipeline.ID}); err != nil {
		db.Instance.Rollback()
		t.Fatal(err)
	}
	db.Instance.Commit()
}
func Test_Clean(t *testing.T) {
	//	db.Instance.Delete(&models.Pipeline{})
	//	db.Instance.Delete(&models.PipelineVersion{})
	//	db.Instance.Delete(&models.Stage{})
	//	db.Instance.Delete(&models.StageVersion{})
	//	db.Instance.Delete(&models.Point{})
	//	db.Instance.Delete(&models.PointVersion{})
}
