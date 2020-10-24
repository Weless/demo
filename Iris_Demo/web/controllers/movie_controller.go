package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"iris_demo/repositories"
	"iris_demo/services"
)

type MovieController struct {
	Ctx iris.Context
}

func (c *MovieController) Get() mvc.View {
	movieRepository := repositories.NewMovieManager()
	movieService := services.NewMovieServiceManger(movieRepository)
	MovieResult := movieService.ShowMovieName()

	return mvc.View{
		Name: "index.html",
		Data: MovieResult,
	}

}
