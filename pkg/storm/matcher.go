package storm

import (
	"reflect"
	"strings"

	"github.com/asdine/storm/q"
)

type like struct {
	val string
}

func Like(fieldName string, val string) q.Matcher {
	return q.NewFieldMatcher(fieldName, &like{val: val})
}

func (l *like) MatchField(v interface{}) (bool, error) {
	refV := reflect.ValueOf(v)
	if refV.Kind() == reflect.String {
		vs := v.(string)
		if strings.Contains(vs, l.val) {
			return true, nil
		}
	}
	return false, nil
}
