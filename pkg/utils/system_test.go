package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGetCpuCores(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want int
// 	}{
// 		{name: RandomString(4), want: 8},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			assert.Equalf(t, tt.want, GetCpuCores(), "GetCpuCores()")
// 		})
// 	}
// }

// func TestGetSN(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want string
// 	}{
// 		{name: RandomString(4), want: "K9GJC09M7W"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			assert.Equalf(t, tt.want, GetSN(), "GetSN()")
// 		})
// 	}
// }

func TestGetUserHomePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: RandomString(4), want: "/Users/hant"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetUserHomePath(), "GetUserHomePath()")
		})
	}
}
