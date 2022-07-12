package website

import (
	"fmt"
	"log"
	"os"
)

// input types in html.
const (
	Button_inputType         = "button"
	Checkbox_inputType       = "checkbox"
	Color_inputType          = "color"
	Date_inputType           = "date"
	Datetime_local_inputType = "datetime-local"
	Email_inputType          = "email"
	File_inputType           = "file"
	Hidden_inputType         = "hidden"
	Image_inputType          = "image"
	Month_inputType          = "month"
	Number_inputType         = "number"
	Password_inputType       = "password"
	Radio_inputType          = "radio"
	Range_inputType          = "range"
	Reset_inputType          = "reset"
	Search_inputType         = "search"
	Submit_inputType         = "submit"
	Tel_inputType            = "tel"
	Text_inputType           = "text"
	Time_inputType           = "time"
	Url_inputType            = "url"
	Week_inputType           = "week"
)

var PageSources = "C://HtmlPages/"

func Html(s string) string {
	return "<html>" + s + "</html>"
}

func Head(s, id, class string) string {
	return fmt.Sprintf(`<head id="%s" class="%s">%s</head>`, id, class, s)
}

func Body(s, id, class string) string {
	return fmt.Sprintf(`<body id="%s" class="%s">%s</body>`, id, class, s)
}

func Title(s, id, class string) string {
	return fmt.Sprintf(`<title id="%s" class="%s">%s</title>`, id, class, s)
}

func Style(s string) string {
	return fmt.Sprintf(`<style>%s</style>`, s)
}

func Form(s, id, class, method, action string) string {
	return fmt.Sprintf(`<form id="%s" class="%s" method="%s" action="%s">%s</form>`, id, class, method, action, s)
}

func Input(id, class, name, typee, value string) string {
	return fmt.Sprintf(`<input id="%s" class="%s" name="%s" type="%s" value="%s">`, id, class, name, typee, value)
}

func Label(s, id, class string) string {
	return fmt.Sprintf(`<label id="%s" class="%s"> %s </label>`, id, class, s)
}

func Heading(s, id, class string, i int) string {
	if ok := 1 <= i && i <= 6; !ok {
		log.Println("Invalid heading index must be 1 <= h[i] <= 6, i is", i)
		return ""
	}
	return fmt.Sprintf(`<h%d id="%s",class="%s">%s</h%d>`, i, id, class, s, i)
}

func Paragraph(s, id, class string) string {
	return fmt.Sprintf(`<p id="%s" class="%s"> %s </p>`, id, class, s)
}

type Page struct {
	Title string
	Body  []byte
}

func NewPage(title string) *Page {
	return &Page{
		Title: title,
		Body:  []byte{},
	}
}

func (p *Page) Write(data []byte) {
	for _, v := range data {
		if string(v) == ">" {
			p.Body = append(p.Body, []byte(string(v)+"\n")...)
		} else {
			p.Body = append(p.Body, v)
		}
	}
}

func (p *Page) Save() error {
	pageFile, err := os.Create(PageSources + p.Title + ".html")
	if err != nil {
		return err
	}
	pageFile.Write(p.Body)
	defer func() {
		err = pageFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func css() string {
	return ` body {
		margin-left:20px ;
		background-color: #ffffff;
	}
	
	  h1 {
		color: blue;
		padding: 20px;
		text-decoration: underline;
	}
	
	  #myForm {
		width: 170px;
		margin-left:  500px;
		margin-bottom: 300px;
	
		background-color: #eeeeee;
	
		padding-left: 5px;
		padding-right: 5px;
		padding-top: 1px;
		border-style:solid ;
	}
	
	.errors {
		border-width: 2px;
		margin-right: 400px;
		margin-left: 400px;
		border-style:solid ;
		border-color:red ;
		color: red;
	}`
}
