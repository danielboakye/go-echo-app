package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danielboakye/go-echo-app/controllers"
	"github.com/danielboakye/go-echo-app/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//  connect to DB
	conn, err := data.OpenDB()
	if err != nil {
		log.Panic("can't connect to postgres")
	}

	// setup config
	app := controllers.Config{
		Repo: data.NewRepository(conn),
	}

	e := app.NewServer()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))

}
