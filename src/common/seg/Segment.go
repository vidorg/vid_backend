package seg

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/huichen/sego"
	"github.com/vidorg/vid_backend/src/config"
	"strings"
)

type SegmentService struct {
	Config *config.ServerConfig `di:"~"`

	Segmenter *sego.Segmenter `di:"-"`
}

func NewSegmentService(dic *xdi.DiContainer) *SegmentService {
	srv := &SegmentService{}
	if !dic.Inject(srv) {
		panic("Inject failed")
	}

	var segmenter sego.Segmenter
	segmenter.LoadDictionary(srv.Config.SearchConfig.DictPath)
	srv.Segmenter = &segmenter

	return srv
}

func (s *SegmentService) Seg(str string) []string {
	segments := s.Segmenter.Segment([]byte(str))
	return sego.SegmentsToSlice(segments, true)
}

func (s *SegmentService) CatAgainst(tokens []string) string {
	sign := "，。、？！；："
	ret := ""
	for _, token := range tokens {
		if !strings.Contains(sign, token) {
			ret += token + " "
		}
	}
	if len(ret) != 0 {
		ret = ret[0 : len(ret)-1]
	}
	return ret
}
