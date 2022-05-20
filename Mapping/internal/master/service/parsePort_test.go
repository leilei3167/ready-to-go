package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsePorts(t *testing.T) {

	happyCase := []struct {
		name string
		args string
		want []int
	}{
		{
			name: "case_01",
			args: "22,24,35-37",
			want: []int{22, 24, 35, 36, 37},
		},
		{
			name: "case_02",
			args: "21,23,21-25",
			want: []int{21, 22, 23, 24, 25},
		},
		{
			name: "case_03",
			args: "21,23,7,5,4,21-25",
			want: []int{4, 5, 7, 21, 22, 23, 24, 25},
		},
	}
	for _, tt := range happyCase {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePorts(tt.args)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)
		})
	}

	failCase := []struct {
		name string
		args string
		want []int
	}{
		{
			name: "case_01",
			args: "dasda,cxz,wq,23",
			want: nil,
		},
		{
			name: "case_02",
			args: "-1,43-23",
			want: nil,
		},
		{
			name: "case_03",
			args: "top100,123",
			want: nil,
		},
	}
	for _, tt := range failCase {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePorts(tt.args)
			require.NotEmpty(t, err)
			require.Equal(t, got, tt.want)
		})

	}

}
