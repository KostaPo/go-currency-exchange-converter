package health

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Health() Response {
	return Response{
		Status: "ok",
	}
}
