package models

const (
	StarPoint  PointKind = "star"
	CheckPoint PointKind = "check"
	EndPoint   PointKind = "end"
)

type PointKind string

type Point struct {
	Id         int64      `json:"id"`
	Pid        int64      `json:"pid"`
	Type       PointKind  `json:"type" binding:"In(star,check,end)"`
	Triggers   string     `json:"triggers"`   // eg : "a,b,c"
	Conditions string     `json:"conditions"` // eg : "a,b,c"
	Status     uint       `json:"Status" gorm:`
	Created    *time.Time `json:"created" `
	Updated    *time.Time `json:"updated"`
	Deleted    *time.Time `json:"deleted"`
}
