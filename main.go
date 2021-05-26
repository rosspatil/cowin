package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type Tmp struct {
	Data []Data
}

type Data struct {
	Name  string
	Date  string
	Age   int
	Total int
	Dose1 int
	Dose2 int
}

var (
	interval = 30 * time.Minute
	pincode  = "450221"
	tmplt    *template.Template
	url      = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByPin?pincode=%s&&date=%s"
)

func init() {
	tmplt = template.Must(template.New("email").Parse(emailTemplate))
}

func main() {
	t := time.NewTicker(interval).C
	getSlot()
	<-t
}

func resolveDate() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	y, m, d := now.Date()
	return fmt.Sprintf("%d-%d-%d", d, m, y)
}

func getSlot() {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(url, pincode, resolveDate()), nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1 Safari/605.1.15")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		shootMail(err.Error())
		return
	}
	ba, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		shootMail(string(ba))
		return
	}
	cowin := Cowin{}
	err = json.Unmarshal(ba, &cowin)
	if err != nil {
		shootMail(err.Error())
		return
	}
	data := Tmp{
		Data: []Data{},
	}
	for i := 0; i < len(cowin.Centers); i++ {
		center := cowin.Centers[i]
		for j := 0; j < len(center.Sessions); j++ {
			session := center.Sessions[j]
			if session.AvailableCapacity > 0 {
				data.Data = append(data.Data, Data{
					Name:  center.Name,
					Date:  session.Date,
					Age:   session.MinAgeLimit,
					Total: session.AvailableCapacity,
					Dose1: session.AvailableCapacityDose1,
					Dose2: session.AvailableCapacityDose2,
				})
			}
		}
	}
	if len(data.Data) == 0 {
		return
	}
	buf := &bytes.Buffer{}
	tmplt.Execute(buf, data)
	shootMail(buf.String())
}

func shootMail(body string) {
	email := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_PASSWORD")
	msg := "From: " + email + "\n" +
		"To: " + email + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"Subject: Vaccination slots details\n\n" + body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", email, password, "smtp.gmail.com"),
		email, []string{email}, []byte(msg))
	if err != nil {
		panic(err)
	}
}

var emailTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <style>
        table {
            font-family: arial, sans-serif;
            border-collapse: collapse;
            width: 100%;
        }

        td,
        th {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }

        tr:nth-child(even) {
            background-color: #dddddd;
        }
    </style>
    <meta charset="UTF-8">
</head>

<body>
    <H4>Vaccination slots are available at the following centers:</H4>
    <table>
        <tr>
            <th>Name</th>
            <th>Date</th>
            <th>Age</th>
            <th>Total</th>
            <th>Dose 1</th>
            <th>Dose 2</th>
        </tr>
        {{range $v := .Data}}
        <tr>
            <td>{{$v.Name}}</td>
            <td>{{$v.Date}}</td>
            <td>{{$v.Age}}</td>
            <td>{{$v.Total}}</td>
            <td>{{$v.Dose1}}</td>
            <td>{{$v.Dose2}}</td>
        </tr>
        {{end}}
    </table>

</body>

</html>`
