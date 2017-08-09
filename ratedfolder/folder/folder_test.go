package folder

import (
	"testing"
	"time"
)

func TestNewFolder(t *testing.T) {
	now := []time.Time{time.Now()}
	tt := []struct {
		name, path string
		panic      bool
		times      []time.Time
	}{
		{name: "okay", panic: false, path: "/foo", times: now},
		{name: "no time value", panic: true, path: "/foo", times: nil},
		{name: "empty path", panic: true, path: "", times: now},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil && tc.panic {
					t.Errorf("%q should not panic", tc.name)
				}
			}()
			New(tc.path, tc.times...)
		})
	}
}
