package controllers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/danielboakye/go-echo-app/data"
	"github.com/labstack/echo"
)

func (app *Config) getAllUsers(c echo.Context) error {
	users, err := app.Repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	return c.JSON(http.StatusOK, users)
}

func (app *Config) getUser(c echo.Context) error {
	id := c.Param("id")

	user, err := app.Repo.GetOne(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

func (app *Config) saveUser(c echo.Context) error {
	var u data.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	eu, err := app.Repo.GetByEmail(u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "processing error"})
	}

	if eu != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "user exists"})
	}

	id, err := app.Repo.Insert(u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "submission failed"})
	}

	u.ID = id
	u.Password = ""

	return c.JSON(http.StatusCreated, u)
}

func (app *Config) updateUser(c echo.Context) error {
	id := c.Param("id")

	var r data.User
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	user, err := app.Repo.GetOne(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	user.Email = r.Email
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Active = r.Active

	err = app.Repo.Update(*user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "update failed"})
	}

	return c.NoContent(http.StatusAccepted)
}

func (app *Config) deleteUser(c echo.Context) error {
	id := c.Param("id")

	err := app.Repo.DeleteByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "bad request"})
	}

	return c.NoContent(http.StatusAccepted)
}
