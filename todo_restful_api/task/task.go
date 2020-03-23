package task

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type status int

const (
	UNKNOWN status = iota
	TODO
	DONE
)

func (s status) String() string {
	switch s {
	case UNKNOWN:
		return "UNKNOWN"
	case TODO:
		return "TODO"
	case DONE:
		return "DONE"
	default:
		return ""
	}
}

func (s status) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "" {
		return nil, errors.New("status.MarshalJSON: unknown value")
	}
	return []byte(fmt.Sprintf("\"%s\"", str)), nil
}

func (s *status) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"UNKNOWN"`:
		*s = UNKNOWN
	case `"TODO"`:
		*s = TODO
	case `"DONE"`:
		*s = DONE
	default:
		return errors.New("status.UnmarshalJSON: unknown value")
	}
	return nil
}

type Deadline struct {
	time.Time
}

func NewDeadline(t time.Time) *Deadline {
	return &Deadline{t}
}

func (d Deadline) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, d.Unix(), 10), nil
}

func (d *Deadline) UnmarshalJSON(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	d.Time = time.Unix(unix, 0)
	return nil
}

type Task struct {
	ID       ID        `json:"id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Status   status    `json:"status,omitempty"`
	Deadline *Deadline `json:"deadline,omitempty"`
	Priority int       `json:"priority,omitempty"`
	SubTasks []Task    `json:"subTasks,omitempty"`
}

func (t Task) String() string {
	check := "v"
	if t.Status != DONE {
		check = " "
	}
	return fmt.Sprintf("[%s] %s %s", check, t.Title, t.Deadline)
}

type List []Task

func (t List) Len() int {
	return len(t)
}

func (t List) Less(i, j int) bool {
	return t[i].ID < t[j].ID
}

func (t List) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type IncludeSubTasks Task

func (t IncludeSubTasks) indentedString(prefix string) string {
	str := prefix + Task(t).String()
	for _, st := range t.SubTasks {
		str += "\n" + IncludeSubTasks(st).indentedString(prefix+"  ")
	}
	return str
}

func (t IncludeSubTasks) String() string {
	return t.indentedString("")
}
