package util

import "reflect"

// 学cloudwego的generic库，受益匪浅
func NewInstance[T any]() T {
	t := TypeOf[T]()
	// 不能直接返回，特殊类型要单独处理
	switch t.Kind() {
	case reflect.Map:
		return reflect.MakeMap(t).Interface().(T)
	case reflect.Slice, reflect.Array:
		return reflect.MakeSlice(t, 0, 0).Interface().(T)
	case reflect.Ptr:
		t = t.Elem()
		origin := reflect.New(t)
		inst := origin
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
			inst = inst.Elem()
			inst.Set(reflect.New(t))
		}

		return origin.Interface().(T)
	default:
		var t T
		return t
	}
}

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
