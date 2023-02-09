package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"temantrip-subscribe/config"
	"temantrip-subscribe/entities"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Message struct {
	Status string `json:"status"`
}

func isEmail(email string) (bool, error) {
	url := "https://isitarealemail.com/api/email/validate?email=" + url.QueryEscape(email)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "bearer "+config.GetConfig().Mail.API)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var message Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return false, err
	}

	if message.Status == "valid" {
		return true, nil
	}

	return false, nil

}

func main() {
	conf := config.GetConfig()
	db := config.InitMySQL(conf)
	config.AutoMigrate(db)

	v := validator.New()
	v.RegisterValidation("alphaspace", func(fl validator.FieldLevel) bool {
		valid := regexp.MustCompile(`^[a-zA-Z\s]+$`)
		return valid.MatchString(fl.Field().String())
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} || ${remote_ip} ${user_agent} || ${method} ${status} ${uri} ${latency_human} ${error}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "STATUS OK")
	})

	e.POST("/register", func(c echo.Context) error {
		req := new(entities.UserRequest)

		err := c.Bind(&req)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(400, echo.Map{
				"message": "something is wrong with request",
			})
		}

		err = v.Struct(req)
		if err != nil {
			messages := echo.Map{}

			for _, field := range err.(validator.ValidationErrors) {
				switch field.Tag() {
				case "required":
					messages[field.Field()] = "this field is required"
				case "lowercase":
					messages[field.Field()] = "this field is must be lowercase"
				case "e164":
					messages[field.Field()] = "this field is must be valid number"
				case "email":
					messages[field.Field()] = "this field is must be valid email"
				case "min":
					messages[field.Field()] = fmt.Sprintf("this field is minimum %s characters", field.Param())
				case "max":
					messages[field.Field()] = fmt.Sprintf("this field is maximum %s characters", field.Param())
				case "alphaspace":
					messages[field.Field()] = "this field is only alpha and space"
				}
			}

			return c.JSON(400, Response{
				Message: "failed create user",
				Errors:  messages,
			})

		}

		row := db.Where("email = ?", req.Email).Find(&entities.User{}).RowsAffected
		if row == 1 {
			return c.JSON(400, Response{
				Message: "failed create user",
				Errors:  "email already exist",
			})
		}

		res, err := isEmail(req.Email)
		if err != nil {
			log.Println(err.Error())
			return c.JSON(400, echo.Map{
				"message": "something is wrong with request",
			})
		}

		if !res {
			return c.JSON(400, Response{
				Message: "failed create user",
				Errors: echo.Map{
					"email": "this field is not found",
				},
			})
		}

		data := entities.User{
			Name:  req.Name,
			Email: req.Email,
			NoHP:  req.NoHP,
		}

		err = db.Create(&data).Error
		if err != nil {
			log.Println(err.Error())
			return c.JSON(500, echo.Map{
				"message": "something is wrong with server",
			})
		}

		return c.JSON(201, echo.Map{
			"message": "success create user",
		})
	})

	e.Logger.Fatal(e.Start(":" + conf.App.APP_PORT))
}
