package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
)

type Person struct {
	Id           int64
	Name         string
	LastName     string
	NationalCode string
}
type Data struct {
	Persons    []Person
	PersonEdit Person
}

var tmpl = template.Must(template.ParseFiles("template/index.html"))
var tmplLogin = template.Must(template.ParseFiles("template/login.html"))
var data = Data{}

func Index(w http.ResponseWriter, r *http.Request) {
	if c, _ := r.Cookie("matin"); c == nil {
		if err := tmplLogin.Execute(w, map[string]interface{}{
			"error": "",
		}); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Println(err)
		}
	}
}
func Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	for index, v := range data.Persons {
		if v.Id == id {
			data.PersonEdit.NationalCode = data.Persons[index].NationalCode
			data.PersonEdit.Name = data.Persons[index].Name
			data.PersonEdit.LastName = data.Persons[index].LastName
			data.PersonEdit.Id = data.Persons[index].Id
			break
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
	}
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		for index, v := range data.Persons {
			if v.Id == id {
				data.Persons[index].NationalCode = r.FormValue("nationalCode")
				data.Persons[index].Name = r.FormValue("name")
				data.Persons[index].LastName = r.FormValue("lastName")
				data.PersonEdit = Person{}
				break
			}
		}
	} else {
		person := Person{}
		person.Id = rand.Int63()
		person.NationalCode = r.FormValue("nationalCode")
		person.Name = r.FormValue("name")
		person.LastName = r.FormValue("lastName")
		data.Persons = append(data.Persons, person)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "admin" && password == "admin" {
		cookie := &http.Cookie{Name: "matin",
			Value:    "maryam",
			MaxAge:   30,
			HttpOnly: true}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		loginJson := map[string]interface{}{
			"error": "Invalid username or password. Please try again",
		}
		if err := tmplLogin.Execute(w, loginJson); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/add", AddPerson)
	http.HandleFunc("/login", Login)
	fmt.Println("Application is Started...")
	http.ListenAndServe(":8080", nil)
}
