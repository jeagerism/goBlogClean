package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jeagerism/goBlogClean/internal/config"
	blogshandlers "github.com/jeagerism/goBlogClean/modules/blogs/blogsHandlers"
	blogsrepositories "github.com/jeagerism/goBlogClean/modules/blogs/blogsRepositories"
	blogsusecases "github.com/jeagerism/goBlogClean/modules/blogs/blogsUsecases"
	middlewareHandler "github.com/jeagerism/goBlogClean/modules/middlewares/mid_hand"
	middlewareUsecase "github.com/jeagerism/goBlogClean/modules/middlewares/mid_use"
	usershandlers "github.com/jeagerism/goBlogClean/modules/users/usersHandlers"
	usersrepositories "github.com/jeagerism/goBlogClean/modules/users/usersRepositories"
	usersusecases "github.com/jeagerism/goBlogClean/modules/users/usersUsecases"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if err := db.PingContext(ctxPing); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	blogsRepositories := blogsrepositories.NewBlogsRepositories(db)
	blogsUseCases := blogsusecases.NewBlogsUsecase(blogsRepositories)
	blogsHandlers := blogshandlers.NewBlogsHandlers(blogsUseCases)

	usersRepositories := usersrepositories.NewUserRepositories(db)
	usersUsecases := usersusecases.NewUsersUsecases(usersRepositories, cfg.JWTSecret)
	usersHandlers := usershandlers.NewUsersHandlers(usersUsecases)

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(splitOrigins(cfg.CORSAllowOrigins), ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: cfg.CORSAllowOrigins != "*",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/healthz"
		},
	}))

	midUse := middlewareUsecase.NewMiddlewareUsecase(cfg.JWTSecret)
	midHand := middlewareHandler.NewMiddlewareHandler(midUse)

	app.Get("/healthz", healthHandler(db))

	app.Get("/", blogsHandlers.FindBlogs)
	app.Get("/:blogId", blogsHandlers.FindBlog)
	app.Post("/post", midHand.CheckToken(), midHand.CheckRole(), blogsHandlers.PostBlog)
	app.Put("/update", midHand.CheckToken(), midHand.CheckRole(), blogsHandlers.UpdateBlog)
	app.Delete("/:blogId", midHand.CheckToken(), midHand.CheckRole(), blogsHandlers.DeleteBlog)

	app.Post("/signup", usersHandlers.Signup)
	app.Post("/login", usersHandlers.Login)

	addr := ":" + cfg.Port

	go func() {
		if err := app.Listen(addr); err != nil {
			log.Printf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("fiber shutdown: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Printf("db close: %v", err)
	}
}

func splitOrigins(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func healthHandler(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "unreachable",
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	}
}
