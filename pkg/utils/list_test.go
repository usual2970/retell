package utils

import "testing"

func TestFindIndex(t *testing.T) {
	type testStruct struct {
		ID   int
		Name string
	}

	type args[K any] struct {
		list     []K
		findFunc func(K) bool
	}
	tests := []struct {
		name string
		args args[testStruct]
		want int
	}{
		{
			name: "find existing item",
			args: args[testStruct]{
				list: []testStruct{
					{ID: 1, Name: "Alice"},
					{ID: 2, Name: "Bob"},
					{ID: 3, Name: "Charlie"},
				},
				findFunc: func(item testStruct) bool {
					return item.ID == 2
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindIndex(tt.args.list, tt.args.findFunc); got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
