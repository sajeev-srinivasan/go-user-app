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

	userResponse := callUsersSource(1)

	var users []dto.User

	for _, user := range userResponse.Users {
		users = append(users, user)
	}
	// Fetching users data concurrently
	channel := make(chan dto.UserResponse, 2)
	wg := sync.WaitGroup{}
	for i := 2; i <= userResponse.TotalPages; i++ {
		wg.Add(1)
		go fetchUsers(i, channel, &wg)
	}
	go func(channel chan dto.UserResponse, wg *sync.WaitGroup) {
		wg.Wait()
		close(channel)
	}(channel, &wg)

	for response := range channel {
		for _, user := range response.Users {
			users = append(users, user)
		}
	}

	return users
}

func callUsersSource(pageNumber int) dto.UserResponse {
	resp, err := http.Get(fmt.Sprintf("https://reqres.in/api/users?page=%d&per_page=2", pageNumber))
	if err != nil {
		return dto.UserResponse{}
	}
	var userResponse dto.UserResponse
	err1 := json.NewDecoder(resp.Body).Decode(&userResponse)
	if err1 != nil {
		println("Decoding error")
	}
	return userResponse
}

func fetchUsers(pageNumber int, channel chan dto.UserResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	userResponse := callUsersSource(pageNumber)
	fmt.Println(userResponse)
	channel <- userResponse
	return
}
