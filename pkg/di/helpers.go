package di

import (
	"github.com/samber/lo"
	"github.com/sarulabs/di/v2"
)

func filterNamesByTag(defs map[string]di.Def, tag string) []string {
	return lo.Keys(lo.PickBy(defs, func(key string, value di.Def) bool {
		return lo.ContainsBy(value.Tags, func(item di.Tag) bool {
			return item.Name == tag
		})
	}))
}
