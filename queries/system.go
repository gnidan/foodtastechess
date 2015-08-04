package queries

type SystemQueries interface {
	GetAnswer(query Query) Answer
}

type SystemQueryService struct {
}

func NewSystemQueryService() *SystemQueryService {
	sqs := new(SystemQueryService)
	return sqs
}

func (s *SystemQueryService) GetAnswer(query Query) Answer {
	return nil
}
