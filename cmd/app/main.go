package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
)

type File struct {
	Name       string `json:"name"`
	IsDir      bool   `json:"isDir"`
	Time       int64  `json:"time"`
	Size       int    `json:"size"`
	ChildCount int    `json:"childCount"`
}

const (
	BasePath = "/app/files/"
)

var g errgroup.Group

func GetBasePath(c echo.Context) string {
	basePath := BasePath
	param1 := c.QueryParam("path")
	if param1 != "" {
		basePath = basePath + strings.ReplaceAll(param1, "}{", "/")
	}
	return basePath
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.GET("/list", listHandler)
	e.GET("/delete", deleteHandler)
	e.POST("/upload", uploadHandler)

	//Static Files
	e.GET("/file", func(c echo.Context) error {
		return c.File(GetBasePath(c))
	})
	e.GET("", func(c echo.Context) error {
		return c.File("./public/index.html")
	})
	e.GET("glyphicons.css", func(c echo.Context) error {
		return c.File("./public/glyphicons.css")
	})
	e.GET("glyphicons-halflings-regular.ttf", func(c echo.Context) error {
		return c.File("./public/glyphicons-halflings-regular.ttf")
	})
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":9055"
	} else {
		port = ":" + port
	}
	mainServer := &http.Server{
		Addr:    port,
		Handler: e,
	}
	g.Go(func() error {
		return mainServer.ListenAndServe()
	})
	log.Println("Service Running at: " + port)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func deleteHandler(c echo.Context) error {
	err := os.Remove(GetBasePath(c))
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "File Deleted Succesfully..!")
}

func listHandler(c echo.Context) error {
	files, err := ioutil.ReadDir(GetBasePath(c))
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var list []File
	if len(files) == 0 {
		return c.String(http.StatusOK, "[]")
	}
	for _, file := range files {
		childsCount := 0
		if file.IsDir() {
			childs, err := ioutil.ReadDir(GetBasePath(c) + "/" + file.Name())
			if err == nil {
				childsCount = len(childs)
			}
		}

		item := File{
			Name:       file.Name(),
			IsDir:      file.IsDir(),
			Size:       int(file.Size()),
			Time:       file.ModTime().UnixMilli(),
			ChildCount: childsCount,
		}
		list = append(list, item)
	}
	return c.JSON(http.StatusOK, list)
}

func uploadHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(GetBasePath(c) + file.Filename)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "File uploaded successfully..!")
}
