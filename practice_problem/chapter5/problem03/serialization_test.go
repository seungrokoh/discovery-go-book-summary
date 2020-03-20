package problem03

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

func TestTask_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		in []byte
		out Task
	} {
		{
			[]byte(`{"Title":"title1","ID":64}`),
			Task {
				Title: "title1",
				ID: 64,
			},
		},
		{
			[]byte(`{"Title":"title2","ID":"64"}`),
			Task{
				Title: "title2",
				ID: 64,
			},
		},
	}


	for i, c := range cases {
		var task Task
		if err := json.Unmarshal(c.in, &task); err != nil {
			log.Println(err)
			return
		}
		if !reflect.DeepEqual(task, c.out) {
			t.Errorf("case index : %d\tgot: %v, want: %v", i, task, c.out)
		}
	}
}
