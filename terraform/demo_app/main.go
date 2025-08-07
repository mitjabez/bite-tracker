package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func dbTest(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return -1, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select count(*) from information_schema.tables")
	if err != nil {
		return -1, fmt.Errorf("Cannot read from DB: %v", err)
	}

	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt, nil
}

func main() {
	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		cnt, err := dbTest(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "Yay! Got %d rows\n", cnt)
	})

	port := 8080
	fmt.Printf("Listening on %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
