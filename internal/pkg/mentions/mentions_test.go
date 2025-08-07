package mentions

import (
	"reflect"
	"testing"

	"github.com/skamenetskiy/messages/internal/entity"
)

func Test_Find(t *testing.T) {
	links := map[string]uint64{
		"@Simon": 1,
		"@Julia": 2,
		"@David": 3,
		"#Simon": 4,
		"#Julia": 5,
		"#David": 6,
		"$Simon": 7,
		"$Julia": 8,
		"$David": 9,
		"~Simon": 10,
		"~Julia": 11,
		"~David": 12,
	}
	type args struct {
		content string
		links   map[string]uint64
	}
	tests := []struct {
		name string
		args args
		want entity.Mentions
	}{
		{"@Simon", args{"Hello @Simon", links}, entity.Mentions{
			{1, 1, [2]uint32{6, 12}},
		}},
		{"@Simon+1", args{"Hello @Simon and @Julia", links}, entity.Mentions{
			{1, 1, [2]uint32{6, 12}},
			{2, 1, [2]uint32{17, 23}},
		}},
		{"@Simon+2", args{"Hello @Simon, @Julia and @David", links}, entity.Mentions{
			{1, 1, [2]uint32{6, 12}},
			{2, 1, [2]uint32{14, 20}},
			{3, 1, [2]uint32{25, 31}},
		}},
		{"@Simon+2 multiline", args{"Hello @Simon\n, @Julia\n and @David\n", links}, entity.Mentions{
			{1, 1, [2]uint32{6, 12}},
			{2, 1, [2]uint32{15, 21}},
			{3, 1, [2]uint32{27, 33}},
		}},
		{"#Simon", args{"Hello #Simon", links}, entity.Mentions{
			{4, 2, [2]uint32{6, 12}},
		}},
		{"#Simon+1", args{"Hello #Simon and #Julia", links}, entity.Mentions{
			{4, 2, [2]uint32{6, 12}},
			{5, 2, [2]uint32{17, 23}},
		}},
		{"#Simon+2", args{"Hello #Simon, #Julia and #David", links}, entity.Mentions{
			{4, 2, [2]uint32{6, 12}},
			{5, 2, [2]uint32{14, 20}},
			{6, 2, [2]uint32{25, 31}},
		}},
		{"#Simon+2 multiline", args{"Hello #Simon\n, #Julia\n and #David\n", links}, entity.Mentions{
			{4, 2, [2]uint32{6, 12}},
			{5, 2, [2]uint32{15, 21}},
			{6, 2, [2]uint32{27, 33}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.content, tt.args.links); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
