package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {

	csvFile, err := os.Open("spotinput.csv")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(csvLines[1][0])

	err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err2 := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err2 != nil {
		log.Fatalf("error connecting to the database: ", err2)
	}
	fmt.Println("Database bzinga opened and ready.")
	defer db.Close()
	_, err4 := db.Exec("truncate table spotinput")
	if err4 != nil {
		log.Fatal("connection problem ", err4)
	}
	indx := 1
	sqlStatement := "insert into spotinput values($1,$2,$3,$4,$5,$6)"
	// sqlStatement := "insert into hourinput values($1,$2,$3,$4,$5)"

	for i := 1; i < len(csvLines); i++ {
		// id, er := strconv.Atoi(csvLines[i][0])
		// if er != nil {
		// 	log.Fatal(er)
		// }
		id := csvLines[i][0]
		cate := csvLines[i][1]
		prod := csvLines[i][2]
		desc := csvLines[i][3]
		mr := csvLines[i][4]
		mr = strings.ReplaceAll(mr, ",", "")
		mr = strings.Trim(mr, " ")
		mrp, er := strconv.Atoi(mr)
		if er != nil {
			log.Fatal(er)
		}
		// mr := csvLines[i][4]
		// bp := csvLines[i][5]

		//used for spot
		_, err3 := db.Exec(sqlStatement, id, cate, prod, desc, mrp, indx)
		// _, err3 := db.Exec(sqlStatement, id, cate, prod, desc, mrp)
		if err3 != nil {
			log.Fatal(err3, ": error in db.exec of sql statement")
		}
		indx = indx + 1
	}

}
