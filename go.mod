module github.com/mirror-media/usersrv

replace github.com/mirror-media/usersrv => /home/baronchiu_mirrormedia_mg/user

go 1.15

require (
	firebase.google.com/go/v4 v4.0.0
	github.com/gin-contrib/static v0.0.0-20200916080430-d45d9a37d28e
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/api v0.32.0
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.2
)
