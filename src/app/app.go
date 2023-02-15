package app

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/event"
	"github.com/we7coreteam/w7-rangine-go/src/components/logger"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/components/validation"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type App struct {
	Name            string
	config          *viper.Viper
	container       container.Container
	providerManager *provider.ProviderManager
	console         *console.Console
}

func NewApp() *App {
	app := &App{
		Name: "rangine",
	}

	app.InitConfig()
	app.InitContainer()
	app.InitProviderManager()
	app.RegisterProviders()

	return app
}

func (app *App) InitConfig() {
	conf := viper.New()
	conf.SetConfigFile("./.env")

	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}

	app.config = conf
}

func (app *App) GetConfig() *viper.Viper {
	return app.config
}

func (app *App) InitContainer() {
	app.container = container.New()
}

func (app *App) GetContainer() container.Container {
	return app.container
}

func (app *App) InitProviderManager() {
	app.providerManager = provider.NewProviderManager(app.container, app.config)
}

func (app *App) GetProviderManager() *provider.ProviderManager {
	return app.providerManager
}

func (app *App) RegisterProviders() {
	app.providerManager.RegisterProvider(new(logger.LoggerProvider)).Register()
	app.providerManager.RegisterProvider(new(event.EventProvider)).Register()
	app.providerManager.RegisterProvider(new(validation.ValidationProvider)).Register()
	app.providerManager.RegisterProvider(new(database.DatabaseProvider)).Register()
	app.providerManager.RegisterProvider(new(redis.RedisProvider)).Register()
}

func (app *App) InitConsole() {
	app.console = new(console.Console)
}

func (app *App) GetConsole() *console.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}

//func (app *App) registerEvent() {
//	app.Event = EventBus.New()
//}
//
//func (app *App) registerValidation() {
//	uni := ut.New(zh.New())
//	lang := app.Config.App.Lang
//	if lang == "" {
//		lang = "zh"
//	}
//
//	app.Translator, _ = uni.GetTranslator(lang)
//	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), app.Translator)
//}
