package go_container

import "reflect"

// Returning "package.struct" from struct instance
func GetPkgPath(s interface{}) string {
	v := reflect.TypeOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.String()
}
