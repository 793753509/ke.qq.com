package model
type CourseInfo struct {
	TableName string `json:"日期,omitempty"`
	ID string `json:"Id"`
	Name string `json:"名称"`
	Teachers string `json:"教师"`
	CourseType string `json:"类别"`
	Price string `json:"价格"`
}
