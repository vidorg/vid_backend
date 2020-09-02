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
func (cmn *CommonService) BuildOrExp(tok string, ids []uint64) string {
	sp := strings.Builder{}
	for _, uid := range ids {
		sp.WriteString(fmt.Sprintf("`%s` = '%d' OR ", tok, uid))
	}
	where := sp.String()[:sp.Len()-4]
	return where
}

type _IdCntPair struct {
	Id  uint64
	Cnt int32
}

type _FromToUidPair struct {
	FromUid uint64
	ToUid   uint64
}

type _UidVidPair struct {
	Uid uint64
	Vid uint64
}
