package introspection

import (
	"sync"

	"gorm.io/gorm/schema"

	"github.com/anthonysyk/versatilego/functional"
)

var columnsCache sync.Map

// This is the testable func, you should normally use GetStructColumns
func getStructColumnsWithCache[T any](model T, cache *sync.Map) []string {
	s, err := schema.Parse(model, cache, schema.NamingStrategy{})
	if err != nil {
		return nil
	}
	return functional.Filter(functional.Map[*schema.Field](s.Fields, func(f *schema.Field) string { return f.DBName }), functional.IsNotEmpty[string])
}

func GetStructColumns[T any](model T) []string {
	return getStructColumnsWithCache(model, &columnsCache)
}
