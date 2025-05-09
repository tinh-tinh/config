package config_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/config/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Namespace(t *testing.T) {
	type Config struct {
		DBHost string `mapstructure:"MYSQL_DBHOST"`
		DBPort string `mapstructure:"MYSQL_DBPORT"`
		DBUser string `mapstructure:"MYSQL_DBUSER"`
		DBPass string `mapstructure:"MYSQL_DBPASS"`
		DBName string `mapstructure:"MYSQL_DBNAME"`
	}

	mysqlController := func(module core.Module) core.Controller {
		ctrl := module.NewController("mysql")

		cfg := config.InjectNamespace[Config](module, "mysql")
		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": cfg.DBPort,
			})
		})

		return ctrl
	}

	mysqlModule := func(module core.Module) core.Module {
		mysql := module.New(core.NewModuleOptions{
			Imports: []core.Modules{
				config.ForFeature[Config]("mysql"),
			},
			Controllers: []core.Controllers{mysqlController},
		})

		return mysql
	}

	mongoController := func(module core.Module) core.Controller {
		ctrl := module.NewController("mongo")

		cfg := config.InjectNamespace[Config](module, "mongo")
		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": cfg.DBPort,
			})
		})

		return ctrl
	}

	mongoModule := func(module core.Module) core.Module {
		mongo := module.New(core.NewModuleOptions{
			Imports: []core.Modules{
				config.ForFeature("mongo", func() *Config {
					return &Config{
						DBHost: os.Getenv("MONGO_DBHOST"),
						DBPort: os.Getenv("MONGO_DBPORT"),
						DBUser: os.Getenv("MONGO_DBUSER"),
						DBPass: os.Getenv("MONGO_DBPASS"),
						DBName: os.Getenv("MONGO_DBNAME"),
					}
				}),
			},
			Controllers: []core.Controllers{mongoController},
		})

		return mongo
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				config.ForRootRaw(".env.example"),
				mongoModule,
				mysqlModule,
			},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/mysql")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, string(`{"data":"3306"}`), string(data))

	resp, err = testClient.Get(testServer.URL + "/api/mongo")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, string(`{"data":"27017"}`), string(data))
}

func Test_NewConfigRaw(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			config.ForRootRaw(),
		},
	})
	require.NotNil(t, appModule)
}

func Test_Nil(t *testing.T) {
	module := core.NewModule(core.NewModuleOptions{})

	cf := config.Inject[Config](module)
	require.Nil(t, cf)

	namespace := config.InjectNamespace[Config](module, "mongo")
	require.Nil(t, namespace)
}
