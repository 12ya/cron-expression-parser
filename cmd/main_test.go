package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		min     int
		max     int
		want    []int
		wantErr bool
	}{
		{
			name: "empty",
			str:  "",
			min:  0,
			max:  59,
			want: nil,
		},
		{
			name: "single",
			str:  "1",
			min:  0,
			max:  59,
			want: []int{1},
		},
		{
			name: "range",
			str:  "1-5",
			min:  0,
			max:  59,
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "range with step",
			str:  "1-5/2",
			min:  0,
			max:  59,
			want: []int{1, 3, 5},
		},
		{
			name: "multiple ranges with step",
			str:  "1-5/2,10-15/3",
			min:  0,
			max:  59,
			want: []int{1, 3, 5, 10, 13},
		},
		{
			name:    "invalid range",
			str:     "1-",
			min:     0,
			max:     59,
			wantErr: true,
		},
		{
			name:    "invalid step",
			str:     "1/2/3",
			min:     0,
			max:     59,
			wantErr: true,
		},
		{
			name:    "invalid step",
			str:     "1-5/2/3",
			min:     0,
			max:     59,
			wantErr: true,
		},
		{
			name:    "invalid step",
			str:     "1-5/2/3,10-15/3",
			min:     0,
			max:     59,
			wantErr: true,
		},
		{
			name:    "invalid range",
			str:     "1-5/2,10-15/3,",
			min:     0,
			max:     59,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.str, tt.min, tt.max)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
