package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uzrnem/go/utils"
)

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//export DESTINATION_DIR=$PWD/dumper/
	fmt.Println("Storing files at : " + utils.ReadEnvOrDefault("DESTINATION_DIR", "/dumper/"))

	e.GET("/:name", func(c echo.Context) error {
		distDir := utils.ReadEnvOrDefault("DESTINATION_DIR", "/dumper/")
		return c.File(distDir + c.Param("name"))
	})
	e.POST("/upload", upload)
	/*
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}))
	*/
	port := utils.ReadEnvOrDefault("SERVER_PORT", "9050")
	e.Logger.Fatal(e.Start(":" + port))
}

func upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	distDir := utils.ReadEnvOrDefault("DESTINATION_DIR", "/dumper/")
	// Destination
	dst, err := os.Create(distDir + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Uploaded successfully :%s", file.Filename))
}
