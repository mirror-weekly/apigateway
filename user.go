package usersrv

type User interface {
	IsSignedIn() (bool, error)
	UpdateInfo(map[string]string) error
	GetInfo() map[string]interface{}
	Delete() error
}

type Service struct {
}

func (s *Service) SignOut(user User) (err error) {
	return err
}

func (s *Service) Update(user User, info map[string]interface{}) (err error) {
	return err
}

func (s *Service) VerifyUser(user User) (ok bool, err error) {
	return ok, err
}
