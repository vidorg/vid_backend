package property

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"log"
	"strings"
)

func (m *PropMapping) ApplyOrderBy(source string) string {
	result := make([]string, 0)
	if source == "" {
		return ""
	}

	sources := strings.Split(source, ",")
	for _, src := range sources {
		src = strings.TrimSpace(src)
		reverse := strings.HasSuffix(src, " desc")
		src = strings.Split(src, " ")[0]

		dest, ok := m.Dict[src]
		if !ok || dest == nil || len(dest.DestProps) == 0 {
			continue
		}

		if dest.Revert {
			reverse = !reverse
		}
		props := dest.DestProps
		props = xslice.ItsOfString(xslice.Reverse(xslice.Sti(props)))

		for _, prop := range props {
			prop += xcondition.IfThenElse(reverse, " DESC", " ASC").(string)
			result = append(result, prop)
		}
	}

	log.Println(result)
	return strings.Join(result, ",")
}
