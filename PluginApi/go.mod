module github.com/BonusPlay/VueHoster/PluginApi

go 1.13

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/sirupsen/logrus v1.4.2
)

require github.com/BonusPlay/VueHoster v0.0.0

replace github.com/BonusPlay/VueHoster v0.0.0 => ../
