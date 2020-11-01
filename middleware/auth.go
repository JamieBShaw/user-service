package middleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const CurrentUserKey = "currentUser"

func Authentication() http.HandlerFunc {

	type UserLoginResponse struct {
		Message string `json:"message"`
		UserId string `json:"user_id"`
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(CurrentUserKey).(int64)
		id := strconv.Itoa(int(userId))

		client := http.Client{
			Timeout: time.Second * 5,
		}

		r, err := http.NewRequest("POST", "http://0.0.0.0:8081/auth/create/"+id, nil)
		if err != nil {
			return
		}

		resp, err := client.Do(r)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		log.Printf("res: %v", string(body))

		response := &UserLoginResponse{
			Message: "User Successfully logged in",
			UserId: id,
		}
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}
