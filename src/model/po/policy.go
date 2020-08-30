package po

type RbacRule struct {
	PType string `gorm:"column:p_type"` //  p  |  g
	V0    string `gorm:"column:v0"`     // sub | sub
	V1    string `gorm:"column:v1"`     // obj | sub
	V2    string `gorm:"column:v2"`     // act |  x
	V3    string `gorm:"column:v3"`
	V4    string `gorm:"column:v4"`
	V5    string `gorm:"column:v5"`
}
