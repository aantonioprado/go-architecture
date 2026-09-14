package usecases

import "github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"

// UserInteractor implements UserInputPort. It depends only on the
// UserRepository port it defines itself; it has no knowledge of whatever
// concrete gateway is plugged in at the composition root.
type UserInteractor struct {
	repo UserRepository
}

func NewUserInteractor(repo UserRepository) UserInputPort {
	return &UserInteractor{
		repo: repo,
	}
}

func (uc *UserInteractor) CreateUser(input CreateUserInput, output UserOutputPort) {
	user, err := entities.NewUser(input.Name, input.Email)
	if err != nil {
		output.PresentError(err)
		return
	}

	if _, err := uc.repo.FindByEmail(user.Email); err == nil {
		output.PresentError(ErrEmailTaken)
		return
	}

	if err := uc.repo.Create(user); err != nil {
		output.PresentError(err)
		return
	}

	output.PresentUserCreated(UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
