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

func (r *RbacRule) ToMap() map[string]interface{} {
	m := map[string]interface{}{"p_type": r.PType}
	if r.V0 != "" {
		m["v0"] = r.V0
	}
	if r.V1 != "" {
		m["v1"] = r.V1
	}
	if r.V2 != "" {
		m["v2"] = r.V2
	}
	if r.V3 != "" {
		m["v3"] = r.V3
	}
	if r.V4 != "" {
		m["v4"] = r.V4
	}
	if r.V5 != "" {
		m["v5"] = r.V5
	}
	return m
}
