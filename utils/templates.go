package utils

import (
	"code.google.com/p/gorilla/sessions"
	"html/template"
	"math/rand"
	"reflect"
)

var FuncMap = template.FuncMap{
	"any": Any,
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

var Store = sessions.NewCookieStore([]byte(randomString(100)))

// any reports whether the first argument is equal to
// any of the remaining arguments.
func Any(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}
