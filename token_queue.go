package gocalc

type queue []*token

func (q *queue) push(t *token) {
	*q = append(*q, t)
}

func (q *queue) pop() *token {
	t := (*q)[0]
	*q = (*q)[1:]
	return t
}

func (q *queue) first() *token {
	return (*q)[0]
}
