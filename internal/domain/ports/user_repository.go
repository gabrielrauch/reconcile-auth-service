package ports

import "github.com/gabrielrauch/reconcile-auth-service/internal/domain/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}
