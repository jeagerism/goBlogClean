package miduse

type middlewareUsecase struct {
	middlewareRepo middlewareRepository.IMiddlewareRepository
}

type IMiddlewareUsecase interface {
}

func NewMiddlewareUsecase()
