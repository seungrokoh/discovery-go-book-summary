package problem04

type Queue struct {
	items []int
}

func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (int, bool) {
	if q.Len() <= 0 {
		return -1, false
	}

	result := q.items[0]
	q.items = q.items[1:]
	return result, true
}

func (q Queue) Len() int {
	return len(q.items)
}