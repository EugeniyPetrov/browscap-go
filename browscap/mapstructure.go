package browscap

import (
	"github.com/go-viper/mapstructure/v2"
	"reflect"
	"strconv"
)

func Int64ToStringHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.Int64 && to.Kind() == reflect.String {
			return strconv.FormatInt(data.(int64), 10), nil
		}

		return data, nil
	}
}
