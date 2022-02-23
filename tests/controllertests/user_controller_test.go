package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"nickname":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			nickname:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"nickname":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "email already taken",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "nickname already taken",
		},
		{
			inputJSON:    `{"nickname":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"nickname": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required nickname",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
		{
			inputJSON:    `{"nickname": "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "required password",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	userSample := []struct {
		id           string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			nickname:   user.Name,
			email:      user.Username,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("cannot convert to json: %v\n", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Name, responseMap["nickname"])
			assert.Equal(t, user.Username, responseMap["username"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var AuthUsername, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("error seeding user: %v\n", err)
	}

	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthUsername = user.Username
		AuthPassword = "password"
	}

	token, err := server.SignIn(AuthUsername, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateNickname string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"nickname":"Grand", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:     200,
			updateNickname: "Grand",
			updateEmail:    "grand@gmail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Woman", "email": "woman@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Man", "email": "man@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Woman", "email": "woman@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Frank", "email": "kenny@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "email already taken",
		},
		{
			// Remember "Kenny Morris" belongs to user 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Kenny Morris", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "nickname already taken",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname":"Kan", "email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "invalid email",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname": "", "email": "kan@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required nickname",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"nickname": "Kan", "email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required email",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"nickname": "Mike", "email": "mike@gmail.com", "password": "password"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "unauthorized",
		},
	}
	// fmt.Println("Auth1: ", reflect.TypeOf(AuthID), "Auth2: ", reflect.TypeOf(AuthID2))
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nickname"], v.updateNickname)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}

		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	var AuthUsername, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("error seeding user: %v\n", err)
	}
	// Get only the first and log him in
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthUsername = user.Username
		AuthPassword = "password"
	}

	token, err := server.SignIn(AuthUsername, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// User 2 trying to use User 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "unauthorized",
		},
	}

	for _, v := range userSample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		req.Header.Set("Authorized", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
