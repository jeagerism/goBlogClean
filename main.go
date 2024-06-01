package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	blogshandlers "github.com/jeagerism/goBlogClean/modules/blogs/blogsHandlers"
	blogsrepositories "github.com/jeagerism/goBlogClean/modules/blogs/blogsRepositories"
	blogsusecases "github.com/jeagerism/goBlogClean/modules/blogs/blogsUsecases"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dataSourceName := "host=localhost user=myuser password=mypassword dbname=mydatabase port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(db)
	}

	blogsRepositories := blogsrepositories.NewBlogsRepositories(db)
	// _ = blogsRepositories

	// blogs, err := blogsRepositories.GetAll()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(blogs)

	// blog, err := blogsRepositories.GetById(1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(blog)

	// post, err := blogsRepositories.Post(mock)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(post)

	blogsUseCases := blogsusecases.NewBlogsUsecase(blogsRepositories)
	blogsHandlers := blogshandlers.NewBlogsHandlers(blogsUseCases)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Get("/", blogsHandlers.FindBlogs)
	app.Listen(":8000")

}
