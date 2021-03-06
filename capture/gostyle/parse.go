package gostyle

import "reflect"

// reflection version

type Info map[string]interface{}

// for interface
func (i Info) Info() interface{} {
	return nil
}

func InfoFromInterface(ptr interface{}, excludes []string) Info {
	rt := reflect.TypeOf(ptr).Elem()
	rv := reflect.ValueOf(ptr).Elem()
	info := Info{}

toplevel:
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)
		for _, name := range excludes {
			if name == rf.Name {
				continue toplevel
			}
		}
		info[rf.Name] = rv.Field(i).Interface()
	}
	return info
}
