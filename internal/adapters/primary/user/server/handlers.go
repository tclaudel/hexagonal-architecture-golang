package server

import "github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"

var (
	_ ServerInterface = (*Handlers)(nil)
)

type Handlers struct {
	user primary.UserUseCase
}

func NewHandlers(userUserCase primary.UserUseCase) *Handlers {
	return &Handlers{
		user: userUserCase,
	}
}
