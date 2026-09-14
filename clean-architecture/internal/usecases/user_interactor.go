package usecases

import "github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"

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

	output.PresentUserCreated(toUserOutput(user))
}

func (uc *UserInteractor) ListUsers(output UserOutputPort) {
	users, err := uc.repo.FindAll()
	if err != nil {
		output.PresentError(err)
		return
	}

	items := make([]UserOutput, 0, len(users))
	for _, user := range users {
		items = append(items, toUserOutput(user))
	}

	output.PresentUserList(ListUsersOutput{Users: items})
}

func (uc *UserInteractor) GetUserById(input GetUserInput, output UserOutputPort) {
	user, err := uc.repo.FindById(input.ID)
	if err != nil {
		output.PresentError(err)
		return
	}

	output.PresentUser(toUserOutput(user))
}

func (uc *UserInteractor) UpdateUser(input UpdateUserInput, output UserOutputPort) {
	existing, err := uc.repo.FindById(input.ID)
	if err != nil {
		output.PresentError(err)
		return
	}

	updated, err := existing.Update(input.Name, input.Email)
	if err != nil {
		output.PresentError(err)
		return
	}

	if updated.Email != existing.Email {
		if _, err := uc.repo.FindByEmail(updated.Email); err == nil {
			output.PresentError(ErrEmailTaken)
			return
		}
	}

	if err := uc.repo.Update(updated); err != nil {
		output.PresentError(err)
		return
	}

	output.PresentUser(toUserOutput(updated))
}

func toUserOutput(user *entities.User) UserOutput {
	return UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
