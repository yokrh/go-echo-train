package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Template for c.Render
type Template struct {
	Templates *template.Template
}

// Render for c.Render
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

// Layout (Common Info)
type Layout struct {
	Title string
}

var layout = Layout{
	"my-title",
}

/* Controller */

func main() {
	// echo
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// template
	t := &Template{
		// Templates: template.Must(template.ParseGlob("templates/index.html")),
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t
	fmt.Println("DEBUG : ", t.Templates.DefinedTemplates())

	// routing
	e.GET("/", handler1)
	e.GET("/hello", handler2)
	e.GET("/hello/:id", handler3)
	e.GET("/hello/user", handler4) // (POST is better)

	e.Logger.Fatal(e.Start(":1323"))
}

/* Helper */

// handler1 (c.String)
func handler1(c echo.Context) error {
	str := getHello()
	str += " (root)"
	return c.String(http.StatusOK, str)
}

// hanlder2 (c.Render, query param)
func handler2(c echo.Context) error {
	n := c.QueryParam("name")
	str := getHelloWithName(n)

	data := struct {
		Layout
		Str string
	}{
		Layout: layout,
		Str:    str,
	}
	return c.Render(http.StatusOK, "index.html", data)
}

// handler3 (c.Render, path param)
func handler3(c echo.Context) error {
	i := c.Param("id")
	str := getHelloWithName(i)

	data := struct {
		Layout
		Str string
	}{
		Layout: layout,
		Str:    str,
	}
	return c.Render(http.StatusOK, "index.html", data)
}

// handler4 (data binding, c.JSON)
func handler4(c echo.Context) error {
	u := new(User)
	// data binding
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest,
			"Request is failed: "+err.Error())
	}
	// validation
	if u.Name == "" || u.Email == "" { // (using "validator" is better)
		const error = "Invalid Params"
		return c.String(http.StatusBadRequest,
			"Request is failed: "+error)
	}

	// str := getHelloWithUser(u)
	return c.JSON(http.StatusOK, u)
}

// User (struct)
type User struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

// function
func getHello() string {
	const str = "Hello, Echo!!!"
	return str
}
func getHelloWithName(name string) string {
	str := "Hello, " + name + "!!!"
	return str
}
func getHelloWithUser(user *User) string {
	str := "Hello, " + user.Name + "(" + user.Email + ")!!!"
	return str
}
