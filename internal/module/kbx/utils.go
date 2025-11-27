package kbx

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

var (
	validKindMap = map[string]reflect.Kind{
		reflect.Struct.String():    reflect.Struct,
		reflect.Map.String():       reflect.Map,
		reflect.Slice.String():     reflect.Slice,
		reflect.Array.String():     reflect.Array,
		reflect.Chan.String():      reflect.Chan,
		reflect.Interface.String(): reflect.Interface,
		reflect.Ptr.String():       reflect.Ptr,
		reflect.String.String():    reflect.String,
		reflect.Int.String():       reflect.Int,
		reflect.Float32.String():   reflect.Float32,
		reflect.Float64.String():   reflect.Float64,
		reflect.Bool.String():      reflect.Bool,
		reflect.Uint.String():      reflect.Uint,
		reflect.Uint8.String():     reflect.Uint8,
		reflect.Uint16.String():    reflect.Uint16,
		reflect.Uint32.String():    reflect.Uint32,
		reflect.Uint64.String():    reflect.Uint64,
	}
)

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetValueOrDefault[T any](value T, defaultValue T) (T, reflect.Type) {
	if !IsObjValid(value) {
		return defaultValue, reflect.TypeFor[T]()
	}
	return value, reflect.TypeFor[T]()
}

func GetValueOrDefaultSimple[T any](value T, defaultValue T) T {
	if !IsObjValid(value) {
		return defaultValue
	}
	return value
}

func IsObjValid(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() == reflect.Ptr && v.Elem().IsNil() {
				return false
			}
			v = v.Elem()
		}
	}
	if _, ok := validKindMap[v.Kind().String()]; !ok {
		return false
	}
	if !v.IsValid() {
		return false
	}
	if v.IsZero() {
		return false
	}
	if v.Kind() == reflect.String && v.Len() == 0 {
		return false
	}
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map || v.Kind() == reflect.Array) && v.Len() == 0 {
		return false
	}
	if v.Kind() == reflect.Bool {
		return true
	}
	return true
}

func IsObjSafe(obj any, strict bool) bool {
	v := reflect.ValueOf(obj)

	// nil pointers or invalid values
	if !v.IsValid() {
		return false
	}
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	// zero value check (different meaning in strict vs resilient mode)
	if v.IsZero() {
		if strict {

			switch v.Kind() {
			case reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64, reflect.String:
				// 0, false, "" são válidos em modo estrito
				return true
			}
		}
		return false
	}

	// empty collections → false no resilient mode
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return !strict
		}
	}

	return true
}

func GetEnvOrDefaultWithType[T any](key string, defaultValue T) T {
	value := os.Getenv(key)
	// Sempre vem texto da env
	if len(value) == 0 {
		return defaultValue
	}
	if reflect.ValueOf(value).CanConvert(reflect.TypeFor[T]()) {
		return reflect.ValueOf(value).Convert(reflect.TypeFor[T]()).Interface().(T)
	}
	var result T
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return defaultValue
	}
	if IsObjValid(result) {
		return result
	}
	return result
}

func GetValueOrDefaultAny[T any](value T, defaultValue T) T {
	if !IsObjValid(value) {
		return defaultValue
	}
	return value
}

func HydrateMapFromEnvOrDefaults[T any](dbType string, target map[string]T, defaults map[string]T, hydrationCtl chan any) map[string]T {
	defer func(hCtl chan any) {
		if r := recover(); r != nil {
			// Handle the panic (e.g., log the error)
			gl.Log("error", fmt.Sprintf("Panic at the Hydration: %v", r))
			if hydrationCtl != nil {
				gl.Log("info", "HydrationCtl", "Async hydration due to panic recovery")
				for key, defaultValue := range defaults {
					target[key] = GetValueOrDefaultAny(target[key], defaultValue)
				}
				hydrationCtl <- r
				return
			}
		}
	}(hydrationCtl)

	for key, defaultValue := range defaults {
		target[key] = GetEnvOrDefaultWithType(dbType+"_"+key,
			GetValueOrDefaultAny(target[key], defaultValue),
		)
	}

	gl.Log("debug", fmt.Sprintf("Hydrated Map for DBType %s: %+v", dbType, target))

	return target
}

func BoolPtr(b bool) *bool {
	return &b
}

func DefaultTrue(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

func DefaultFalse(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func PtrInt64(n int64) *int64 {
	return &n
}

func PtrInt(n int64) *int {
	ni := int(n)
	return &ni
}

func PtrDuration(d time.Duration) *time.Duration {
	if d == 0 {
		return nil
	}
	return &d
}
