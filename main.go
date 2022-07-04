package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func css() string {
	return `body {
	margin-left:20px ;
	background-color: #ffffff;
	}
	
	h1 {
	color: blue;
	padding: 20px;
	text-decoration: underline;
	}
	
	div {
	width: 170px;
	margin-left:  500px;
	margin-bottom: 300px;
	height: 387px;
	
	background-color: #eeeeee;
	
	padding-left: 5px;
	padding-right: 5px;
	padding-bottom: 5px;
	padding-top: 1px;
	border-width: 0px;
	border-style:solid ;
	border-radius: 25px;
	border-color:black ;
	
	}
	
	.errors {
	border-width: 2px;
	  border-style:solid ;
	  border-color:red ;
	  color: red;
	}`
}

const dataBasePath string = "C://UsersData/"
const port string = ":8080"

var (
	errInvalidPassword              = errors.New("the password is not as the verified one?")
	errNoRegisteredOrDeletedAccount = errors.New("no account registered with this gmail? Make one?")
	errLostOrDeletedData            = errors.New("lost or deletd data? Visit your actions log?")
	errAlreadyOccupiedGmail         = errors.New("an account already registered with the given gmail? User another one or check your gmail links?")
)

type Name string

type User struct {
	Firstname       Name `json:"firstname"`
	Lastname        Name `json:"lastname"`
	Nation          Name `json:"nation"`
	Birthday        `json:"birthday"`
	Gendre          `json:"gendre"`
	VIPSubscription `json:"VIPSubscription"`
	Gmail           `json:"gmail"`
	Password        `json:"password"`
	ID              int `json:"id"`
}

type Day uint // 1 <= Day <= 31

type Month time.Month // 1 <= Month <= 12

type Year uint

type Gmail string

type Gendre string

type Password string

type Birthday struct {
	Day
	Month
	Year
}

// Unused
type VIPSubscription struct {
	lasts        time.Duration
	subscription bool
}

var (
	firstName      Name
	lastName       Name
	nation         Name
	gendre         Gendre
	gmail          Gmail
	password       Password
	verifyPassword Password
)

func signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/signup/" {
		switch r.Method {

		case "GET":
			http.ServeFile(w, r, "signUp.html")
		case "POST":

			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			// here we get the signup information.

			firstName = Name(r.FormValue("firstName"))
			lastName = Name(r.FormValue("lastName"))
			nation = Name(r.FormValue("nation"))
			gendre = Gendre(r.FormValue("gendre"))
			gmail = Gmail(r.FormValue("gmail"))
			password = Password(r.FormValue("password"))
			verifyPassword = Password(r.FormValue("verifyPassword"))

			// signUp.html but with filled values
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

<form method="POST" action="/">     
  <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
  
  <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>

  <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>

  <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>

  <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
  
  <p><label> Password:</label><input name="password" type="password" value="%s" />              </p>

  <p><label> Verify password:</label><input name="verifyPassword" type="password" value="%s" /> </p>

  <input type="submit" value="submit" />

</form>
<p class ="errors">Error: %s</p>
</div>
</body>
</html>
			,`, css(), gmail, gendre, firstName, lastName, nation, password, verifyPassword, errInvalidPassword.Error())
				return
			}

			var sliceOfErrors = []error{}
			sliceOfErrors = append(sliceOfErrors, firstName.Check(), lastName.Check(), nation.Check(), gendre.Check(), gmail.Check(), password.Check(), verifyPassword.Check())

			// we check for errors...
			for _, v := range sliceOfErrors {
				if v != nil {
					// signUp.html but with filled values.
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
   
  <form method="POST" action="/">     
      <p><label> Gmail address:</label><input name="gmail" type="email" value="%s" />               </p>
      
      <p><label> Gendre</label><input name="gendre" type="text" value="%s" />                       </p>

      <p><label> First Name:</label><input name="firstName" type="text" value="%s" />               </p>

      <p><label> Last Name:</label><input name="lastName" type="text" value="%s" />                 </p>

      <p><label> Nation:</label><input name="nation" type="text" value="%s" />                      </p>
      
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

			user := NewUser(firstName, lastName, nation, gendre, gmail, password, 22, 7, 2006)
			path := string(gmail)

			// we make a new folder at C://UsersData/(entered gmail)dir to store the json info.
			// like C://UsersData/example@gmail.comdir

			err := os.Mkdir(dataBasePath+path+"dir", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// like C://UsersData/example@gmail.comdir/data.json .
			jsonFile, err := os.Create(dataBasePath + path + "dir/" + "data.json")
			if err != nil {
				http.Error(w, errAlreadyOccupiedGmail.Error(), 1)
				return
			}

			// writes to data.json the content of the sign up.
			err = json.NewEncoder(jsonFile).Encode(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// close data.json and if it fails the whole server falls down.
			err = jsonFile.Close()
			if err != nil {
				log.Fatal(err)
			}

			// writes signUpAccepted.html to the user screen so he knows his account has been created.
			http.ServeFile(w, r, "signUpAccepted.html")

		default:
			fmt.Fprintf(w, "Only GET and POST methods are supported.")
		}
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/login/" {
		switch r.Method {

		case "GET":
			http.ServeFile(w, r, "logIn.html")
		case "POST":

			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			// here we get the login information.

			gmail = Gmail(r.FormValue("gmail"))
			password = Password(r.FormValue("password"))

			// we check for errors...
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
				
					<form method="POST" action="/">
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

			// open the data.json file from "C://UsersData/example@gmail.comdir"
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
				
					<form method="POST" action="/">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
						<input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, errNoRegisteredOrDeletedAccount.Error())
				return
			}

			// read the content of data.json and put them in user
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
				
					<form method="POST" action="/">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
						<input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, errLostOrDeletedData.Error())
				return
			}

			// compare the json password and the form password
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
				
					<form method="POST" action="/">
						<p><label> Gmail:</label><input name="gmail" type="email" value="%s" />                       </p>
				
						<p><label> Password:</label><input name="password" type="password" value="%s" />              </p>
				
						<input type="submit" value="log in" />
					</form>
					<p class ="errors">Error: %s</p>
				</body>
				
				</html>`, css(), gmail, password, err.Error())
				return
			}

			// close data.json and if it fails the whole server falls down.
			err = jsonFile.Close()
			if err != nil {
				log.Fatal(err)
			}

			// writes logInAccepted.html to the user screen so he knows he logedin.
			http.ServeFile(w, r, "logInAccepted.html")
		default:
			fmt.Fprintf(w, "Only GET and POST methods are supported.")
		}
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	http.ServeFile(w, r, "welcome.html")
}

func main() {
	http.HandleFunc("/welcome/", welcome)
	http.HandleFunc("/signup/", signUp)
	http.HandleFunc("/login/", logIn)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// Check() checks g contains at the end @gmail.com.
func (g Gmail) Check() (err error) {
	exists := strings.Contains(string(g), "@gmail.com")
	if !exists {
		return errors.New("invalid gmail? Must contains '@gmail.com'! ")
	}
	return
}

// Check() checks if 1 <= d <= 31.
func (d Day) Check() (err error) {
	if 1 <= d {
		if d <= 31 {
			return
		} else {
			return fmt.Errorf("Day > 31. %d", d)
		}
	} else {
		return fmt.Errorf("Day < 1. %d", d)
	}
}

// Check() checks if 1 <= m <= 12.
func (m Month) Check() (err error) {
	if lessThan1, greaterThan12 := m < 1, m > 12; lessThan1 || greaterThan12 {
		return errors.New("invalid month of birth? Must be 1 <= m <= 12! ")
	}
	return
}

// Check() checks if 1 <= y.
func (y Year) Check() (err error) {
	if lessThan1 := y < 1; lessThan1 {
		return errors.New("invalid year of birth? Must be 1 < y! ")
	}
	return
}

// Check() checks if d, m and y are valid.
func (b Birthday) Check() (err error) {
	if ok := b.Day.Check() != nil; ok {
		return b.Day.Check()
	} else if ok := b.Month.Check() != nil; ok {
		return b.Month.Check()
	} else if ok := b.Year.Check() != nil; ok {
		return b.Year.Check()
	} else {
		return
	}
}

// Check() checks if the gendre is male or female.
func (g Gendre) Check() (err error) {
	s := strings.ToLower(string(g))
	if s != "male" && s != "female" {
		return errors.New("invalid Gendre?")
	}
	return
}

// Check() checks if the name is empty.
func (n Name) Check() (err error) {
	if n == "" {
		return errors.New("empty name?")
	}
	return
}

// Check() checks if the password is empty or contains illegal charaters.
func (p Password) Check() (err error) {
	if p == "" {
		return errors.New("empty password?")
	}

	if strings.Contains(string(p), "(") || strings.Contains(string(p), ")") || strings.Contains(string(p), "{") || strings.Contains(string(p), "}") || strings.Contains(string(p), "[") || strings.Contains(string(p), "]") || strings.Contains(string(p), "#") {
		return errors.New("illegal Charater(s)? Must not contains (){}[]#?")
	}
	return
}

func (p Password) Compare(p2 Password) (err error) {
	if ok := p != p2; ok {
		return errors.New("password is not correct?")
	}
	return nil
}

var (
	name Name
	pass Password
	g    Gmail
	gen  Gendre
	b    Birthday
	sub  VIPSubscription
)

// NewUser() returns a pointer to a user with the given parameters.
func NewUser(firstname, lastname, nation Name, gendre Gendre, gmail Gmail, password Password, day, month, year int) *User {
	return &User{
		Firstname: firstname,
		Lastname:  lastname,
		Nation:    nation,
		Birthday: Birthday{
			Day:   Day(day),
			Month: Month(month),
			Year:  Year(year),
		},
		Gendre:          gendre,
		VIPSubscription: sub,
		Gmail:           gmail,
		Password:        password,
	}
}
