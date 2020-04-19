package util

import (
    "reflect"
)

func Contain(target interface{}, obj interface{}) bool {
    targetValue := reflect.ValueOf(target)
    switch reflect.TypeOf(target).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if targetValue.Index(i).Interface() == obj {
                return true
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            return true
        }
    }

    return false
}

func SliceRemoveUint32(target *[]uint32, obj uint32) {
    s := *target
    p := -1
    for i, v := range s {
        if v == obj {
            p = i
            break
        }
    }
    if p != -1 {
      *target = append(s[:p], s[p+1:]...)
    }
}

func SliceRemoveInt(target *[]int, obj int) {
    s := *target
    p := -1
    for i, v := range s {
        if v == obj {
            p = i
            break
        }
    }
    if p != -1 {
      *target = append(s[:p], s[p+1:]...)
    }
}

func SliceRemoveString(target *[]string, obj string) {
    s := *target
    p := -1
    for i, v := range s {
        if v == obj {
            p = i
            break
        }
    }
    if p != -1 {
      *target = append(s[:p], s[p+1:]...)
    }
}
