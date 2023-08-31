package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	repo "github.com/A1sal/AvitoTech/repository"
	sg "github.com/A1sal/AvitoTech/segment"
	u "github.com/A1sal/AvitoTech/user"
	usseg "github.com/A1sal/AvitoTech/usersegment"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

func initRepo(conn ksql.DB) {
	var dbSegment *sg.SegmentActualDatabase = sg.NewSegmentActualDatabase(conn)
	var dbUserSegment *usseg.UserSegmentActualDatabase = usseg.NewUserSegmentActualDatabase(conn)
	var dbUser *u.UserActualDatabase = u.NewUserActualDatabase(conn)
	serviceRepo = *repo.NewServiceMockRepository(dbSegment, dbUserSegment, dbUser)
}

func initDatabaseConnection(ctx context.Context) (ksql.DB, error) {
	psPort := os.Getenv("POSTGRES_PORT")
	psDb := os.Getenv("POSTGRES_DB")
	psUser := os.Getenv("POSTGRES_USER")
	psPassword := os.Getenv("POSTGRES_PASSWORD")
	connString := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s", psUser, psPassword, psPort, psDb)
	conn, err := kpgx.New(ctx, connString, ksql.Config{})
	if err != nil {
		fmt.Printf("Could not connect to postgres: %v\n", err)
		return conn, err
	}
	return conn, nil
}

var serviceRepo repo.ServiceMockRepository

func main() {
	ctx := context.Background()
	godotenv.Load()
	conn, err := initDatabaseConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.Mkdir("data", os.ModeDir); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	defer conn.Close()
	initRepo(conn)
	port := flag.String("port", "3003", "port the server will listen to")
	flag.Parse()
	r := chi.NewRouter()
	r.Get("/", helloRootHandler)
	r.Post("/{segmentName}", createSegmentHandler)
	r.Delete("/{segmentName}", deleteSegmentHandler)
	r.Put("/modify-user-segments", modifyUserSegments)
	r.Get("/{userId:[0-9]+}", getUserSegments)
	r.Get("/{userId:[0-9]+}/{year:[0-9]+}/{month:[0-9]+}", getUserSegmentsInPeriod)
	r.Get("/user-report/{filename}", downloadUserReport)

	log.Fatal(http.ListenAndServe(":"+*port, r))
}
