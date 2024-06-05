package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	blogshandlers "github.com/jeagerism/goBlogClean/modules/blogs/blogsHandlers"
	blogsrepositories "github.com/jeagerism/goBlogClean/modules/blogs/blogsRepositories"
	blogsusecases "github.com/jeagerism/goBlogClean/modules/blogs/blogsUsecases"
	usershandlers "github.com/jeagerism/goBlogClean/modules/users/usersHandlers"
	usersrepositories "github.com/jeagerism/goBlogClean/modules/users/usersRepositories"
	usersusecases "github.com/jeagerism/goBlogClean/modules/users/usersUsecases"
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

	//==>BLOG ZONE
	blogsRepositories := blogsrepositories.NewBlogsRepositories(db)
	blogsUseCases := blogsusecases.NewBlogsUsecase(blogsRepositories)
	blogsHandlers := blogshandlers.NewBlogsHandlers(blogsUseCases)
	//==>BLOG ZONE

	//==> USER ZONE
	usersRepositories := usersrepositories.NewUserRepositories(db)
	usersUsecases := usersusecases.NewUsersUsecases(usersRepositories)
	usersHandlers := usershandlers.NewUsersHandlers(usersUsecases)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/", blogsHandlers.FindBlogs)
	app.Get("/:blogId", blogsHandlers.FindBlog)
	app.Post("/post", blogsHandlers.PostBlog)
	app.Put("/update", blogsHandlers.UpdateBlog)
	app.Delete("/:blogId", blogsHandlers.DeleteBlog)

	app.Post("/signup", usersHandlers.Signup)
	app.Post("/login", usersHandlers.Login)
	app.Listen(":8000")

}
