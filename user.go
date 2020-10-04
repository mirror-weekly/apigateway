package usersrv

type User interface {
	IsSignedIn() (bool, error)
	UpdateInfo(map[string]string) error
	GetInfo() map[string]string
	Delete() error
}
