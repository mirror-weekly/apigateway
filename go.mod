module github.com/mirror-media/usersrv

replace github.com/mirror-media/usersrv => /home/baronchiu_mirrormedia_mg/user

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0 // indirect
	cloud.google.com/go/storage v1.12.0 // indirect
	firebase.google.com/go v3.13.0+incompatible // indirect
	github.com/gin-contrib/static v0.0.0-20200916080430-d45d9a37d28e
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.7.0
	gopkg.in/yaml.v2 v2.3.0
)
