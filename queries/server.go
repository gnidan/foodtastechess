package queries

type SystemQueries interface {
	GetAnswer(query Query) Answer
}
