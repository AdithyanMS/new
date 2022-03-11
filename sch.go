package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func schedule(db *sql.DB) {
	endSpot := time.Now()
	startSpot := time.Now()
	var durationSpot int
	updateSpot := time.Now()
	var intervalSpot int
	var count int
	var id int
	var ref string
	var category string
	var prod_name string
	var desc string
	var mrp int
	var indx int
	tid := 216
	tcount := 40
	out1 := "insert into spottable values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)"

	rowLine := db.QueryRow("select * from logic where id=$1", 2)
	err1 := rowLine.Scan(&id, &count, &startSpot, &durationSpot, &intervalSpot, &endSpot)
	if err1 != nil {
		log.Fatal(err1)
	}
	spotIndex := 1
	spotAuction := 0
	_, err4 := db.Exec("truncate table spottable")
	if err4 != nil {
		log.Fatal("connection problem ", err4)
	}
	fmt.Println("spottable truncated")
	for !startSpot.After(endSpot) {
		rows, err5 := db.Query("select * from spotinput where indx<=$1 order by random()", tid)
		if err5 != nil {
			log.Fatal(err5)
		}
		for rows.Next() {
			fmt.Println("Spotindex: ", spotIndex)
			fmt.Printf("%s\n", startSpot)
			err := rows.Scan(&ref, &category, &prod_name, &desc, &mrp, &indx)
			if err != nil {
				log.Fatal(err)
			}
			if startSpot.After(endSpot) {
				break
			}
			if startSpot.Hour() >= 23 {
				fmt.Println("Inside dark hour")
				tid = tid + tcount
				nt := startSpot.String()
				ntl := strings.Split(nt, " ")
				ntl[1] = "09:00:00"
				ntl = ntl[:2]
				nt = ntl[0] + "T" + ntl[1] + "+00:00"
				start2, _ := time.Parse("2006-01-02T15:04:05Z07:00", nt)
				start2 = start2.AddDate(0, 0, 1)
				startSpot = start2
				break
			}
			updateSpot = startSpot.Add(time.Duration(durationSpot) * time.Minute)
			base := 0.99 * float64(mrp)
			_, execerr := db.Exec(out1, spotIndex, spotAuction, startSpot, updateSpot, 1, 1, ref, category, prod_name, desc, mrp, 1, int(base))
			if execerr != nil {
				fmt.Println(execerr)
			}
			if spotIndex%count == 0 {
				spotAuction++
				startSpot = startSpot.Add(time.Duration(intervalSpot) * time.Minute)
			}
			spotIndex++
		}
	}
}
func main() {
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal(err1)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	fmt.Println("Database student opened and ready.")
	schedule(db)
}
