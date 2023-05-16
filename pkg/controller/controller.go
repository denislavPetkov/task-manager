package controller

import (
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"go.uber.org/zap"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	taskdb "github.com/denislavpetkov/task-manager/pkg/database/task"
	userdb "github.com/denislavpetkov/task-manager/pkg/database/user"
	"github.com/denislavpetkov/task-manager/pkg/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type Controller interface {
	Start() error
}

type controller struct {
	ginRouter *gin.Engine

	userDb userdb.Redis
	taskDb taskdb.Mongodb
}

func NewController() Controller {
	return &controller{}
}

const (
	serverPort = ":8081"
)

var (
	logger *zap.Logger
)

func (c *controller) Start() error {
	logger = zap.L().Named("controller")

	logger.Info("Starting controller")

	err := c.init()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize controller, error: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Starting server on port %s", serverPort))

	// Start the server
	err = c.ginRouter.Run(serverPort)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start web server, error: %v", err))
		return err
	}

	return nil
}

func (c *controller) init() error {
	c.initGin()

	err := c.initRedis()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize redis database, error: %v", err))
		return err
	}

	err = c.initMongodb()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize mongodb database, error: %v", err))
		return err
	}

	c.registerPaths()

	logger.Info("Controller initialization successful")

	return nil
}

func (c *controller) initGin() {
	logger.Info("Initializating gin")

	gin.SetMode(gin.ReleaseMode)

	c.ginRouter = gin.Default()

	c.ginRouter.SetTrustedProxies(nil)
	c.ginRouter.Static("/static/css", "./static/css")
	c.ginRouter.LoadHTMLGlob("static/*.html")
}

func (c *controller) initRedis() error {
	logger.Info("Initializating redis database")

	redisConfig, err := userdb.NewRedisConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get redis config, error: %v", err))
		return err
	}

	c.userDb = userdb.NewRedis(redisConfig)

	err = c.initSessionStore(redisConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to initialize session store, error: %v", err))
	}

	return err
}

// the gin engine has to be initialised before
func (c *controller) initSessionStore(redisConfig userdb.RedisConfig) error {
	logger.Info("Initializating session store")

	redisStoreToken, err := crypto.GenerateToken(redisConfig.GetPassword())
	if err != nil {
		return err
	}

	store, err := redis.NewStore(10, "tcp", redisConfig.GetAddress(), redisConfig.GetPassword(), []byte(redisStoreToken))
	if err != nil {
		return err
	}

	store.Options(sessions.Options{
		MaxAge:   constants.SessionCookieTtl,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	csrfSecretToken, err := crypto.GenerateToken(redisConfig.GetPassword())
	if err != nil {
		return err
	}

	c.ginRouter.Use(sessions.Sessions(constants.SessionStoreName, store))
	c.ginRouter.Use(csrf.Middleware(csrf.Options{
		Secret: csrfSecretToken,
		ErrorFunc: func(c *gin.Context) {
			c.String(401, "CSRF token mismatch")
			c.Abort()
		},
	}))

	return nil
}

func (c *controller) initMongodb() error {
	logger.Info("Initializating mongodb database")

	mongodbConfig, err := taskdb.NewMongodbConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get mongodb config, error: %v", err))
		return err
	}

	c.taskDb, err = taskdb.NewMongodb(mongodbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get mongodb instance, error: %v", err))
	}

	return err
}

func (c *controller) registerPaths() {
	logger.Info("Registering web paths")

	c.ginRouter.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	c.ginRouter.POST("/audio", c.postAudio)

	c.ginRouter.GET("/register", c.getRegister)
	c.ginRouter.POST("/register", c.postRegister)

	c.ginRouter.GET("/login", c.getLogin)
	c.ginRouter.POST("/login", c.postLogin)

	authenticated := c.ginRouter.Group("/")
	authenticated.Use(middleware.Authentication())
	{

		authenticated.GET("/passChange", c.getPasswordChange)
		authenticated.POST("/passChange", c.postPasswordChange)

		authenticated.GET("/tasks/edit/:title", c.getUpdateTask)
		authenticated.POST("/tasks/edit/:title", c.postUpdateTask)
		authenticated.DELETE("/tasks/delete/:title", c.deleteTask)
		authenticated.POST("/tasks/completed/:title", c.postCompleteTask)
		authenticated.GET("/tasks", c.getTasks)

		authenticated.GET("/tasks/new", c.getNewTask)
		authenticated.POST("/tasks/new", c.postNewTask)

		authenticated.GET("/logout", c.getLogout)

	}
}
