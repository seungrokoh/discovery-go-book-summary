package problem03

import (
	"encoding/json"
	"errors"
	"strconv"
)

type ID int64
type Task struct {
	Title string `json:"title"`
	ID    ID     `json:"id"`
}

func (i *ID) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	//fmt.Println(reflect.TypeOf(v))
	switch x := v.(type) {
	case string:
		s, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
		*i = ID(s)
	case float64:
		*i = ID(x)
	default:
		return errors.New("ID.UnmarshalJSON: unknown value")
	}
	return nil
}
