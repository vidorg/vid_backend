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

// Merge 2 `[]*_IdCntPair` slices to `[]*[2]int32`, using bucket.
func (cmn *CommonService) MergeIdCntPairs(pairs *[2][]*_IdCntPair, uids []uint64) []*[2]int32 {
	// bucket
	bucket := make(map[uint64][2]int32, len(uids))
	for _, subing := range pairs[0] {
		bucket[subing.Id] = [2]int32{subing.Cnt, 0}
	}
	for _, suber := range pairs[1] {
		if arr, ok := bucket[suber.Id]; !ok {
			bucket[suber.Id] = [2]int32{0, suber.Cnt}
		} else {
			bucket[suber.Id] = [2]int32{arr[0], suber.Cnt}
		}
	}

	// out
	out := make([]*[2]int32, len(uids))
	for idx, uid := range uids {
		arr, ok := bucket[uid]
		if ok {
			out[idx] = &arr
		} else {
			out[idx] = &[2]int32{}
		}
	}

	return out
}

// Merge 2 `[]*_FromToUidPair` slices to `[]*[2]bool`, using bucket.
func (cmn *CommonService) MergeFromToUidPairs(pairs *[2][]*_FromToUidPair, uids []uint64) []*[2]bool {
	// bucket
	bucket := make(map[uint64][2]bool, len(uids))
	for _, subing := range pairs[0] {
		bucket[subing.ToUid] = [2]bool{true, false}
	}
	for _, suber := range pairs[1] {
		if arr, ok := bucket[suber.FromUid]; !ok {
			bucket[suber.FromUid] = [2]bool{false, true}
		} else {
			bucket[suber.FromUid] = [2]bool{arr[0], true}
		}
	}

	// out
	out := make([]*[2]bool, len(uids))
	for idx, uid := range uids {
		arr, ok := bucket[uid]
		if ok {
			out[idx] = &arr
		} else {
			out[idx] = &[2]bool{}
		}
	}

	return out
}

type _IdCntPair struct {
	Id  uint64
	Cnt int32
}

type _FromToUidPair struct {
	FromUid uint64
	ToUid   uint64
}
