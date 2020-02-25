package property

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
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
		for i, j := 0, len(props)-1; i < j; i, j = i+1, j-1 {
			props[i], props[j] = props[j], props[i] // reverse
		}

		for _, prop := range props {
			prop += xcondition.IfThenElse(reverse, " DESC", " ASC").(string)
			result = append(result, prop)
		}
	}

	log.Println(result)
	return strings.Join(result, ",")
}
