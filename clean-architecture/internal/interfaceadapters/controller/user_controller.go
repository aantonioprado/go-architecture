package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/presenter"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

// UserController depends on usecases.UserInputPort, the interface, not on
// the concrete UserInteractor.
type UserController struct {
	usecase usecases.UserInputPort
}

func NewUserController(usecase usecases.UserInputPort) *UserController {
	return &UserController{
		usecase: usecase,
	}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	out := presenter.NewHTTPUserPresenter(w)

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		out.PresentError(err)
		return
	}

	c.usecase.CreateUser(usecases.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}, out)
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	out := presenter.NewHTTPUserPresenter(w)

	c.usecase.ListUsers(out)
}

func (c *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	out := presenter.NewHTTPUserPresenter(w)

	c.usecase.GetUserById(usecases.GetUserInput{
		ID: chi.URLParam(r, "id"),
	}, out)
}
