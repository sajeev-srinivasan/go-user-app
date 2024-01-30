package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"user-app/internal/app/dto"
)

type Users interface {
	GetUsers() []dto.User
}

type users struct{}

func NewUsers() Users {
	return &users{}
}

func (u users) GetUsers() []dto.User {

	resp, err := http.Get("https://reqres.in/api/users?page=1&per_page=2")
	if err != nil {
		return nil
	}
	var userResponse dto.UserResponse
	err1 := json.NewDecoder(resp.Body).Decode(&userResponse)
	if err1 != nil {
		println("Decoding error")
	}

	var users []dto.User

	for _, user := range userResponse.Users {
		users = append(users, user)
	}

	channel := make(chan dto.UserResponse, 5)
	fmt.Println("Before fetching users concurrently")
	// Fetching users data concurrently
	wg := sync.WaitGroup{}
	for i := 2; i <= userResponse.TotalPages; i++ {
		fmt.Println("inside loop")
		wg.Add(1)
		go fetchUsers(i, channel, &wg)
	}
	wg.Wait()
	close(channel)
	fmt.Println("After fetching users concurrently")

	for response := range channel {
		for _, user := range response.Users {
			users = append(users, user)
		}
	}

	fmt.Println("End of GetUsers")
	return users
}

func fetchUsers(pageNumber int, channel chan dto.UserResponse, wg *sync.WaitGroup) {
	fmt.Println("Inside fetch users")
	resp, err := http.Get(fmt.Sprintf("https://reqres.in/api/users?page=%d&per_page=2", pageNumber))
	if err != nil {
		return
	}

	var userResponse dto.UserResponse
	err1 := json.NewDecoder(resp.Body).Decode(&userResponse)
	if err1 != nil {
		println("Decoding error")
	}

	fmt.Println(userResponse)

	channel <- userResponse
	fmt.Println("About to be done with fetch users")
	wg.Done()
	fmt.Println("Done with fetch users")
}
