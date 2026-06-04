package api

import (
	"testing"

	"github.com/spf13/viper"
)

func TestAdminTokenValidityInMinutes(t *testing.T) {
	const key = "internal.admin-token-validity-minutes"
	// Ensure global viper state is restored after the test.
	t.Cleanup(func() { viper.Set(key, nil) })

	cases := []struct {
		name     string
		set      bool
		value    int64
		expected int64
	}{
		{name: "unset falls back to default", set: false, expected: DefaultAdminTokenValidityInMinutes},
		{name: "zero falls back to default", set: true, value: 0, expected: DefaultAdminTokenValidityInMinutes},
		{name: "negative falls back to default", set: true, value: -5, expected: DefaultAdminTokenValidityInMinutes},
		{name: "positive overrides default", set: true, value: 30, expected: 30},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.set {
				viper.Set(key, tc.value)
			} else {
				viper.Set(key, nil)
			}
			if got := adminTokenValidityInMinutes(); got != tc.expected {
				t.Fatalf("adminTokenValidityInMinutes() = %d, want %d", got, tc.expected)
			}
		})
	}
}
