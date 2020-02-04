package seg

import (
	"fmt"
	"github.com/huichen/sego"
	"testing"
)

func TestSegmentService_Seg(t *testing.T) {
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("F:/Projects/vid/vid_backend/src/common/seg/dictionary.txt")

	str := "支持普通和搜索引擎两种分词模式，支持用户词典、词性标注，可运行JSON RPC服务。"
	fmt.Println(sego.SegmentsToString(segmenter.Segment([]byte(str)), true))
	fmt.Println(sego.SegmentsToString(segmenter.Segment([]byte(str)), false))
	fmt.Println(sego.SegmentsToSlice(segmenter.Segment([]byte(str)), true))
	fmt.Println(sego.SegmentsToSlice(segmenter.Segment([]byte(str)), false))
}
