package main

type User struct {
	Id             int    `gorm:"column:id"`
	Title          string `gorm:"column:title"`
	Author         string `gorm:"column:author"`
	CaseContent    string `gorm:"column:case_content"`
	ProductContent string `gorm:"column:product_content"`
	ClassifyId     int    `gorm:"column:classify_id"`
	Banner         string `gorm:"column:banner"`
	Order          int    `gorm:"column:order"`

	Active bool `gorm:"column:active"`
	//TimeCreate string `gorm:"column:time_create"`
	//TimeUpdate string `gorm:"column:time_update"`
}

func (User) TableName() string {
	return "academic_case"
}

func (User) SchemaName() string {
	return "crodis"
}
