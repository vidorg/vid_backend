package service

import (
	"fmt"
	"strings"
)

type CommonService struct{}

func NewCommonService() *CommonService {
	return &CommonService{}
}

// Build `xxx = yyy OR ...` where expression.
func (cmn *CommonService) BuildOrExpr(tok string, ids []uint64) string {
	if len(ids) == 0 {
		return ""
	}

	sp := strings.Builder{}
	for _, uid := range ids {
		sp.WriteString(fmt.Sprintf("`%s` = '%d' OR ", tok, uid))
	}
	where := sp.String()[:sp.Len()-4]
	return where
}

type _IdScanResult struct {
	Id uint64
}

type _IdCntScanResult struct {
	Id  uint64
	Cnt int32
}
