package surrealgo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Headers struct {
	Username string
	Password string
	NS       string
	DB       string
	Accept   string
}

type Session struct {
	URL     string
	Headers *Headers
	Client  *http.Client
}

type Response struct {
	Time   string                   `json:"time"`
	Status string                   `json:"status"`
	Detail string                   `json:"detail"`
	Result []map[string]interface{} `json:"result"`
}

func CreateSession(
	URL string,
	user string,
	pass string,
	NS string,
	DB string,
	Accept string,
) *Session {
	session := &Session{
		URL,
		&Headers{
			user,
			pass,
			NS,
			DB,
			Accept,
		},
		&http.Client{},
	}

	return session
}

func addHeaders(request *http.Request, headers *Headers) *http.Request {
	request.Header.Set("NS", headers.NS)
	request.Header.Set("DB", headers.DB)
	request.Header.Set("Accept", headers.Accept)
	request.SetBasicAuth(headers.Username, headers.Password)

	return request
}

func (session Session) Sql(query string) *Response {
	var content *bytes.Buffer = bytes.NewBuffer([]byte(query))

	req, err := http.NewRequest("POST", (session.URL + "/sql"), content)
	if err != nil {
		log.Fatal(err)
	}

	addHeaders(req, session.Headers)

	resp, err := session.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := []Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	return &response[0]
}
