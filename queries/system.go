package queries

type SystemQueries interface {
	AnswerQuery(query Query) interface{}
}

type SystemQueryService struct {
}

func NewSystemQueryService() *SystemQueryService {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) AnswerQuery(query Query) interface{} {
	return nil
}
