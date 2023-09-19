package services

import "github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"

type UserUseCase struct {
	userRepository secondary.UserRepository
}

func NewUserUseCase(userRepository secondary.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}
