package dto

import (
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("RbacRuleDto", "rbac rule response").
			Properties(
				goapidoc.NewProperty("p_type", "string", true, "rule type"),
				goapidoc.NewProperty("v0", "string", true, "v0"),
				goapidoc.NewProperty("v1", "string", true, "v1"),
				goapidoc.NewProperty("v2", "string", true, "v2"),
				goapidoc.NewProperty("v3", "string", true, "v3"),
				goapidoc.NewProperty("v4", "string", true, "v4"),
				goapidoc.NewProperty("v5", "string", true, "v5"),
			),
	)
}

type RbacRuleDto struct {
	PType string `json:"p_type"` //  p  |  g
	V0    string `json:"v_0"`    // sub | sub
	V1    string `json:"v_1"`    // obj | sub
	V2    string `json:"v_2"`    // act |  x
	V3    string `json:"v_3"`
	V4    string `json:"v_4"`
	V5    string `json:"v_5"`
}

func BuildRbacRuleDto(rule *po.RbacRule) *RbacRuleDto {
	if rule == nil {
		return nil
	}
	return &RbacRuleDto{
		PType: rule.PType,
		V0:    rule.V0,
		V1:    rule.V1,
		V2:    rule.V2,
		V3:    rule.V3,
		V4:    rule.V4,
		V5:    rule.V5,
	}
}

func BuildRbacRuleDtos(rules []*po.RbacRule) []*RbacRuleDto {
	out := make([]*RbacRuleDto, len(rules))
	for idx, rule := range rules {
		out[idx] = BuildRbacRuleDto(rule)
	}
	return out
}
