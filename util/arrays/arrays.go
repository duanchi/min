package arrays

import "reflect"

func ContainsString (array []string, needle string) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}


func ContainsInt (array []int64, needle int64) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsFloat (array []float64, needle float64) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}

func ContainsStruct (array []interface{}, needle interface{}) (index int, has bool) {
	has = false
	index = -1

	for i, spec := range array {
		if reflect.DeepEqual(spec, needle) {
			has = true
			index = i
			return
		}
	}
	return
}
