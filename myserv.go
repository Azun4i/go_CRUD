package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	//"text/template"
	"html/template"
)

type Citi struct {
	Id   int
	Name string
	City string
}

func dbConnect() (db *sql.DB) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "usergo"
		password = "0000"
		dbname   = "mydbgo"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Errorf("can't open db ", err)
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()

	// потом удалить
	//mans := []Citi{}
	//rows, err := db.Query("select * from city")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//for rows.Next() {
	//	p := Citi{}
	//	err := rows.Scan(&p.Name, &p.City)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	mans = append(mans, p)
	//}
	//for _, p := range mans {
	//	fmt.Println(p.Id, p.Name, p.City)
	//}

	selDB, err := db.Query("SELECT * FROM city ORDER BY id DESC")
	if err != nil {
		log.Fatal("can't SELECT in db ", err)
	}
	defer selDB.Close()

	pg := []Citi{}
	for selDB.Next() {
		p := Citi{}
		err := selDB.Scan(&p.Id, &p.Name, &p.City)
		if err != nil {
			panic(err.Error())
		}
		pg = append(pg, p)
	}

	//err = tmpl.ExecuteTemplate(w, "Index", pg)
	err = tmpl.ExecuteTemplate(w, "index.html", pg)
	if err != nil {
	}
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()

	numId := r.URL.Query().Get("Id")
	selectDB, err := db.Query("SELECT * FROM city WHERE id=$1", numId)
	if err != nil {
		fmt.Errorf("can't show ", err, numId)
	}
	p := Citi{}
	for selectDB.Next() {
		err := selectDB.Scan(&p.Id, &p.Name, &p.City)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	//tmpl.ExecuteTemplate(w, "Show", p)
	tmpl.ExecuteTemplate(w, "show.html", p)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	//tmpl.ExecuteTemplate(w, "New", nil)
	tmpl.ExecuteTemplate(w, "new.html", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()

	numId := r.URL.Query().Get("Id")
	selectDB, err := db.Query("SELECT name,city,id FROM city WHERE id=$1", numId)
	if err != nil {
		fmt.Errorf("can't show ", err, numId)
	}
	p := Citi{}
	for selectDB.Next() {
		err := selectDB.Scan(&p.Name, &p.City, &p.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	tmpl.ExecuteTemplate(w, "edit.html", p)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		//insertDB, err := db.Prepare("INSERT INTO city (name,city) VALUES ($1,$2)")
		//if err != nil {
		//	fmt.Errorf("can't Insert ", err, name, city)
		//}
		_, err := db.Exec("INSERT INTO city (name, city) VALUES ($1,$2)", name, city)
		if err != nil {
			fmt.Errorf("can't Insert ", err, name, city)
		}
		log.Println("INSERT new:  name: " + name + " | city: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("id")
		//insForm, err := db.Prepare("UPDATE city SET name=$1, city=$2 WHERE id=$3")
		//if err != nil {
		//	panic(err.Error())
		//}
		//_, err = insForm.Exec(name, city, id)
		//if err != nil {
		//	panic(err.Error())
		//}
		_, err := db.Exec("UPDATE city SET name=$1, city=$2 WHERE id=$3", name, city, id)
		if err != nil {
		}
		log.Println("Update: id: " + id + "| name: " + name + " | city: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	emp := r.URL.Query().Get("Id")
	delForm, err := db.Prepare("DELETE FROM city WHERE id=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("Delete item with id: " + emp)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}

//migrate create -ext sql -dir migrations create_city для миграции
//Database connection strings are specified via URLs. The URL format is driver dependent but generally has the form: dbdriver:
//username:password@host:port/dbname?param1=true&param2=false
//migrate -path migrations -database "postgres://usergo:0000@localhost:5432/mydbgo?sslmode=disable" up
