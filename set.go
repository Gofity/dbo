package dbo

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/Gofity/gokit"
	"golang.org/x/exp/constraints"
)

type xSetConstraint interface {
	constraints.Integer | constraints.Float | ~string
}

type Set[T xSetConstraint] struct {
	Data gokit.Array[T]
}

func (x *Set[T]) Scan(value any) (err error) {
	switch data := value.(type) {
	case []T:
		x.Data = data

	case gokit.Array[T]:
		x.Data = data

	default:
		var entry T

		x.Data = gokit.Array[T]{}
		vtype := reflect.TypeOf(entry)

		gokit.SplitFn(value, ",", func(data gokit.String) {
			var err error

			switch vtype.Kind() {
			case reflect.Pointer:
				entry = reflect.New(vtype.Elem()).Interface().(T)
				err = json.Unmarshal([]byte(data), entry)
			default:
				entry = reflect.Zero(vtype).Interface().(T)
				err = json.Unmarshal([]byte(data), &entry)
			}

			if err == nil {
				x.Data.Append(entry)
			}
		})
	}

	return
}

func (x Set[T]) Value() (value driver.Value, err error) {
	value = x.Data.Join(",")
	return
}

func (x *Set[T]) UnmarshalJSON(b []byte) (err error) {
	value := gokit.Array[T]{}

	if err = json.Unmarshal(b, &value); err != nil {
		return
	}

	return x.Scan(value)
}

func (x Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Data)
}
