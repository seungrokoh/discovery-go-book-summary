package problem04

type Task struct {
	Title string	`json:"title,omitempty"`
	Priority int 	`json:"priority,omitempty"`
}

type TaskList []Task

func (t TaskList) Len() int {
	return len(t)
}

func (t TaskList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Priority가 낮을 수록 우선순위가 높은 것
func (t TaskList) Less(i, j int) bool {
	return t[i].Priority < t[j].Priority
}