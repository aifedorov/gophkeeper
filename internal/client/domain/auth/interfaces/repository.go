package interfaces

//go:generate mockgen -source=interfaces.go -destination=mock_repository_test.go -package=auth

type Repository interface {
	Save(session Session) error
	Load() (Session, error)
	Delete() error
}
