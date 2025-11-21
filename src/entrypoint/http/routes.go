package http

import (
	"github.com/anaclaraddias/feedmaker-students/src/entrypoint/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HandlerInterface interface {
	Handle(c *gin.Context)
}

const (
	HealthCheck = "healthCheck"

	Login      = "login"
	CreateUser = "createUser"

	CreateFeedback       = "createFeedback"
	ListStudentFeedbacks = "listStudentFeedbacks"
	ListTeacherFeedbacks = "listTeacherFeedbacks"
	FindFeedback         = "findFeedback"
)

type Routes struct {
	*gin.Engine
	Handlers map[string]HandlerInterface
	db       *gorm.DB
}

func NewRoutes(
	app *gin.Engine,
	db *gorm.DB,
) *Routes {
	return &Routes{
		app,
		createMapOfHandlers(db),
		db,
	}
}

func (routes *Routes) Register() {
	routes.GET(
		"/health",
		routes.Handlers[HealthCheck].Handle,
	)

	routes.POST(
		"/login",
		routes.Handlers[Login].Handle,
	)

	routes.POST(
		"/user",
		routes.Handlers[CreateUser].Handle,
	)

	routes.POST(
		"/feedback",
		routes.Handlers[CreateFeedback].Handle,
	)

	routes.GET(
		"/student/:id/feedbacks",
		routes.Handlers[ListStudentFeedbacks].Handle,
	)

	routes.GET(
		"/teacher/:id/feedbacks",
		routes.Handlers[ListTeacherFeedbacks].Handle,
	)

	routes.GET(
		"/feedback/:id",
		routes.Handlers[FindFeedback].Handle,
	)
}

func createMapOfHandlers(db *gorm.DB) map[string]HandlerInterface {
	return map[string]HandlerInterface{
		HealthCheck:          handler.NewHealthHandler(db),
		Login:                handler.NewLoginHandler(db),
		CreateUser:           handler.NewCreateUserHandler(db),
		CreateFeedback:       handler.NewCreateFeedbackHandler(db),
		ListStudentFeedbacks: handler.NewListStudentFeedbacksHandler(db),
		ListTeacherFeedbacks: handler.NewListTeacherFeedbacksHandler(db),
		FindFeedback:         handler.NewFindFeedbackHandler(db),
	}
}
