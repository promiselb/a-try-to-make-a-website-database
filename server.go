package website

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewServer() *Server {
	srv := &Server{
		Router: mux.NewRouter().StrictSlash(true),
		datas:  Slice[Data]{},
	}
	srv.Routes()
	return srv
}

func (srv *Server) Routes() {
	(srv.HandleFunc("/signup", srv.SignUp())).Methods("GET", "POST")
	(srv.HandleFunc("/login", srv.LogIn())).Methods("GET", "POST")
	(srv.HandleFunc("/login/{gmail}", srv.Profile())).Methods("GET")
	(srv.HandleFunc("/", srv.Welcome())).Methods("GET")
}

func (srv *Server) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			http.ServeFile(w, r, htmlPagesPath+"signUp.html")
		case "POST":

			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			firstName = Name(r.FormValue("firstName"))
			lastName = Name(r.FormValue("lastName"))
			nation = Name(r.FormValue("nation"))
			gendre = Gendre(r.FormValue("gendre"))
			gmail = Gmail(r.FormValue("gmail"))
			password = Password(r.FormValue("password"))
			sbirthday := r.FormValue("birthday")
			verifyPassword = Password(r.FormValue("verifyPassword"))

			birthday, err = Sbirthday(sbirthday)

			if password != verifyPassword {
				fmt.Fprintf(w, `
				<!DOCTYPE html>
	<html>
	
	<head>
	<title>Sign Up</title>
	<style>
	%s
	</style>
	
	</head>
	
	<body>
	<h1>Sign up for free</h1>
	<div>
	
	<form id="myForm" method="POST" action="/signup">     
	  <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
	  
	  <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>
	
	  <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>
	
	  <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>
	
	  <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
	  
	  <p><label> Birthday:</label><input name="birthday" type="date" value="%s" />                       </p>

	  <p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
	
	  <p><label> Verify password:</label><input name="verifyPassword" type="password" value="%s" /> </p>
	
	  <input type="submit" value="submit" />
	
	</form>
	<p class ="errors">Error: %s</p>
	</div>
	</body>
	</html>
				,`, css(), gmail, gendre, firstName, lastName, nation, sbirthday, password, verifyPassword, errInvalidPassword.Error())
				return
			}

			var sliceOfErrors = []error{}
			sliceOfErrors = append(sliceOfErrors, firstName.Check(), lastName.Check(), nation.Check(), gendre.Check(), gmail.Check(), password.Check(), verifyPassword.Check(), err, birthday.Check())

			for _, v := range sliceOfErrors {
				if v != nil {
					fmt.Fprintf(w, `
					<!DOCTYPE html>
	<html>
	
	<head>
	  <title>Sign Up</title>
	  <style>
	  %s
	</style>
	
	</head>
	
	<body>
	<h1>Sign up for free</h1>
	  <div>
	   
	  <form id="myForm" method="POST" action="/signup">     
		  <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
		  
		  <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>
	
		  <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>
	
		  <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>
	
		  <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
		  
		  <p><label> Birthday:</label><input name="birthday" type="date" value="" />                    </p>

		  <p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
	
		  <p><label> Verify password:</label><input name="verifyPassword" type="password" value="%s" /> </p>
	
		  <input type="submit" value="submit" />
	
	  </form>
	  <p class ="errors">Error: %s</p>
	</div>
	</body>
	</html>
					,`, css(), gmail, gendre, firstName, lastName, nation, password, verifyPassword, v.Error())
					return
				}
			}

			user := NewUser(firstName, lastName, nation, gendre, gmail, password, int(birthday.Day), int(birthday.Month), int(birthday.Year), uuid.New())
			path := string(gmail)

			err = os.Mkdir(dataBasePath+path+"dir", os.ModePerm)
			if err != nil {

				fmt.Fprintf(w, `
					<!DOCTYPE html>
	<html>
	
	<head>
	  <title>Sign Up</title>
	  <style>
	  %s
	</style>
	
	</head>
	
	<body>
	<h1>Sign up for free</h1>
	  <div>
	   
	  <form id="myForm" method="POST" action="/signup">     
		  <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
		  
		  <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>
	
		  <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>
	
		  <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>
	
		  <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
		  
		  <p><label> Birthday:</label><input name="birthday" type="date" value="" />                    </p>

		  <p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
	
		  <p><label> Verify password:</label><input name="verifyPassword" type="password" value="%s" /> </p>
	
		  <input type="submit" value="submit" />
	
	  </form>
	  <p class ="errors">Error: %s</p>
	</div>
	</body>
	</html>
					,`, css(), gmail, gendre, firstName, lastName, nation, password, verifyPassword, errAlreadyOccupiedGmail.Error())
				return
			}

			jsonFile, err := os.Create(dataBasePath + path + "dir/" + "data.json")
			if err != nil {
				fmt.Fprintf(w, `
					<!DOCTYPE html>
	<html>
	
	<head>
	  <title>Sign Up</title>
	  <style>
	  %s
	</style>
	
	</head>
	
	<body>
	<h1>Sign up for free</h1>
	  <div>
	   
	  <form id="myForm" method="POST" action="/signup">     
		  <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
		  
		  <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>
	
		  <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>
	
		  <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>
	
		  <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
		  
		  <p><label> Birthday:</label><input name="birthday" type="date" value="" />                    </p>

		  <p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
	
		  <p><label> Verify password:</label><input name="verifyPassword" type="password" value="%s" /> </p>
	
		  <input type="submit" value="submit" />
	
	  </form>
	  <p class ="errors">Error: %s</p>
	</div>
	</body>
	</html>
					,`, css(), gmail, gendre, firstName, lastName, nation, password, verifyPassword, errAlreadyOccupiedGmail.Error())
				return
			}

			err = json.NewEncoder(jsonFile).Encode(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = jsonFile.Close()
			if err != nil {
				log.Fatal(err)
			}

			http.ServeFile(w, r, "C://Go/src/golang-book/1.18/website/html/signUpAccepted.html")

		default:
			fmt.Fprintf(w, "Only GET and POST methods are supported.")
		}
	}
}

func (srv *Server) LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			http.ServeFile(w, r, htmlPagesPath+"logIn.html")
		case "POST":

			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			gmail = Gmail(r.FormValue("gmail"))
			password = Password(r.FormValue("password"))

			for _, err := range []error{gmail.Check(), password.Check()} {
				if err != nil {
					fmt.Fprintf(w, `<!DOCTYPE html>
				<html>
				
				<head>
					<title>Sign Up</title>
				<style>
				%s
					</style>
				</head>
				<body>
					<h1>Log in</h1>
				
					<form id="myForm" method="POST" action="/login">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
				      <input type="submit" value="log in" /> 
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, err.Error())
					return
				}
			}
			user := &User{}
			path := string(gmail)

			jsonFile, err := os.Open(dataBasePath + path + "dir/" + "data.json")
			if err != nil {
				fmt.Fprintf(w, `<!DOCTYPE html>
				<html>
				
				<head>
					<title>Sign Up</title>
				<style>
				%s
					</style>
				</head>
				<body>
					<h1>Log in</h1>
				
					<form id="myForm"id="myForm" method="POST" action="/login">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
					   <input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, errNoRegisteredOrDeletedAccount.Error())
				return
			}

			err = json.NewDecoder(jsonFile).Decode(user)
			if err != nil {
				fmt.Fprintf(w, `<!DOCTYPE html>
				<html>
				
				<head>
					<title>Sign Up</title>
				<style>
				%s
					</style>
				</head>
				<body>
					<h1>Log in</h1>
				
					<form id="myForm" method="POST" action="/login">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
						<input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, errLostOrDeletedData.Error())
				return
			}

			err = password.Compare(user.Password)
			if err != nil {
				fmt.Fprintf(w, `<!DOCTYPE html>
				<html>
				
				<head>
					<title>Sign Up</title>
				<style>
				%s
				</style>
				</head>
				<body>
					<h1>Log in</h1>
				
					<form id="myForm" method="POST" action="/login">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
						 <input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, err.Error())
				return
			}

			err = jsonFile.Close()
			if err != nil {
				log.Fatal(err)
			}
			redirectURL := fmt.Sprintf(`http://localhost%s/login/%s`, Port, gmail)
			fmt.Println(redirectURL)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)

		default:
			fmt.Fprintf(w, "Only GET and POST methods are supported.")
		}
	}
}

func (srv *Server) Profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["gmail"]
		jsonFile, err := os.Open(dataBasePath + key + "dir/data.json")
		if err != nil {
			fmt.Fprint(w, Html(Head(Title("Profile", "", "")+Style(css()), "", "")+Body(Heading("Your Profile infromation:", "", "", 1)+Paragraph(err.Error(), "", ""), "", "")))
			return
		}
		defer func() {
			err = jsonFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		stats, err := jsonFile.Stat()
		if err != nil {
			fmt.Fprint(w, Html(Head(Title("Profile", "", "")+Style(css()), "", "")+Body(Heading("Your Profile infromation:", "", "", 1)+Paragraph(err.Error(), "", ""), "", "")))
			return
		}

		b := stats.Size()
		data := make([]byte, b)
		_, err = jsonFile.Read(data)
		fmt.Fprint(w, Html(Head(Title("Profile", "", "")+Style(css()), "", "")+Body(Heading("Your Profile infromation:", "", "", 1)+Paragraph(string(data), "", ""), "", "")))
	}
}

func (srv *Server) Welcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, htmlPagesPath+"welcome.html")
	}
}
