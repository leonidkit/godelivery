package main

import (
	"godelivery/internal/api"
	"godelivery/internal/converter/xml2json"
	"godelivery/internal/storage/redis"
	"godelivery/pkg/logger/logrus.go"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type GoDelivery struct {
	apiServer *api.APIServer
}

func main() {

	lgr := logrus.New()

	err := godotenv.Load(".env.example")
	if err != nil {
		lgr.Fatal(".env file parse error: " + err.Error())
	}

	rc, err := redis.New(os.Getenv("REDIS_URL"), time.Second*1800)
	if err != nil {
		lgr.Fatalf("redis connection error: %v", err)
	}
	cnv := xml2json.New()

	apiSrv := api.New(cnv, rc, lgr)

	app := GoDelivery{
		apiServer: apiSrv,
	}

	lgr.Fatal(app.apiServer.Run(os.Getenv("HOST"), os.Getenv("PORT")))
}
