package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coderconvoy/dbase"
	"github.com/coderconvoy/gojs"
	"github.com/coderconvoy/siteman/usr"
)

func LoginHandler(uu []usr.Usr, sc *dbase.SessionControl) MuxFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Check for match
		found := -1
		for k, v := range uu {
			if v.Username == r.FormValue("username") && v.Password.Check(r.FormValue("password")) {
				found = k
				break
			}
		}
		if found == -1 {
			//TODO Add Message somewhere
			http.Redirect(w, r, "/", 303)
			return
		}
		sc.Login(w, uu[found])
		//Add to sessioncontrol
		//point to home
		http.Redirect(w, r, "/home", 303)

	}
}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.Write(IndexPage().Bytes())
}

func main() {
	usrn := flag.Bool("usr", false, "Create or Edit a User")
	usrf := flag.String("usrf", "usrdata.json", "Set Userdata file")

	flag.Parse()

	if *usrn {
		usr.RunUserFunc(*usrf)
		return
	}

	users, err := usr.LoadUsers(*usrf)
	if err != nil {
		fmt.Println(err)
		return
	}

	sesh := dbase.NewSessionControl(time.Minute * 15)

	gojs.Single.AddFuncs(Asset, AssetDir)

	http.HandleFunc("/save", NewHandler(users, sesh, FileSaver))
	http.HandleFunc("/delete", NewHandler(users, sesh, FileDeleter))
	http.HandleFunc("/move", NewHandler(users, sesh, FileMover))
	http.HandleFunc("/upload", NewHandler(users, sesh, FileUploader))
	http.HandleFunc("/ass/", gojs.AssetHandler("/ass/", gojs.Single))
	http.HandleFunc("/home", NewHandler(users, sesh, HomeView))
	http.HandleFunc("/login", LoginHandler(users, sesh))
	http.HandleFunc("/usr/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	log.Fatal(http.ListenAndServe(":8090", nil))

}
