package controller

import (
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/crypto"

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

func (c *controller) Start() error {
	err := c.init()
	if err != nil {
		return err
	}

	c.ginRouter.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	c.ginRouter.GET("/register", c.getRegister)
	c.ginRouter.POST("/register", c.postRegister)

	c.ginRouter.GET("/login", c.getLogin)
	c.ginRouter.POST("/login", c.postLogin)

	authenticated := c.ginRouter.Group("/")
	authenticated.Use(middleware.Authentication())
	{

		authenticated.GET("/tasks/edit/:title", c.getUpdateTask)
		authenticated.POST("/tasks/edit/:title", c.postUpdateTask)
		authenticated.DELETE("/tasks/delete/:title", c.deleteTask)
		authenticated.POST("/tasks/completed/:title", c.postCompleteTask)
		authenticated.GET("/tasks", c.getTasks)

		authenticated.GET("/tasks/new", c.getNewTask)
		authenticated.POST("/tasks/new", c.postNewTask)

		authenticated.GET("/logout", c.getLogout)

	}

	// Start the server
	return c.ginRouter.Run(":8081")
}

func (c *controller) init() error {
	c.initGin()

	err := c.initRedis()
	if err != nil {
		return err
	}

	err = c.initMongodb()
	if err != nil {
		return err
	}

	return nil
}

func (c *controller) initRedis() error {
	redisConfig, err := userdb.NewRedisConfig()
	if err != nil {
		return err
	}

	c.userDb = userdb.NewRedis(redisConfig)

	return c.initSessionStore(redisConfig)
}

// the gin engine has to be initialised before
func (c *controller) initSessionStore(redisConfig userdb.RedisConfig) error {
	store, err := redis.NewStore(10, "tcp", redisConfig.GetAddress(), redisConfig.GetPassword(), []byte(crypto.GenerateToken(redisConfig.GetPassword())))
	if err != nil {
		return err
	}

	store.Options(sessions.Options{
		MaxAge:   constants.SessionCookieTtl,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	c.ginRouter.Use(sessions.Sessions(constants.SessionStoreName, store))
	c.ginRouter.Use(csrf.Middleware(csrf.Options{
		Secret: crypto.GenerateToken(redisConfig.GetPassword()),
		ErrorFunc: func(c *gin.Context) {
			c.String(401, "CSRF token mismatch")
			c.Abort()
		},
	}))

	return nil
}

func (c *controller) initMongodb() error {
	mongodbConfig, err := taskdb.NewMongodbConfig()
	if err != nil {
		return err
	}

	c.taskDb, err = taskdb.NewMongodb(mongodbConfig)

	return err
}

func (c *controller) initGin() {
	c.ginRouter = gin.Default()

	c.ginRouter.SetTrustedProxies(nil)
	c.ginRouter.Static("/static/css", "./static/css")
	c.ginRouter.LoadHTMLGlob("static/*.html")
}
