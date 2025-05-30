package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	totalUsers = 10_000_000
	batchSize  = 5000
	minFriends = 1
	maxFriends = 20
)

var cities = []string{
	"Москва", "Лондон", "Нью-Йорк", "Берлин", "Париж", "Пекин",
	"Токио", "Сидней", "Дубай", "Рим",
}

func main() {
	dsn := "DB_url"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	// 1. Вставка пользователей
	log.Println("Generating users...")
	usernames := make([]string, 0, totalUsers)
	for i := 0; i < totalUsers; i++ {
		usernames = append(usernames, fmt.Sprintf("user%d", i+1))
	}
	// Перемешать, если хочется случайности распределения по friends
	rand.Shuffle(len(usernames), func(i, j int) {
		usernames[i], usernames[j] = usernames[j], usernames[i]
	})

	for i := 0; i < totalUsers; i += batchSize {
		var sb strings.Builder
		sb.WriteString("INSERT INTO users (username, city, points) VALUES ")
		args := make([]interface{}, 0, batchSize*3)
		for j := 0; j < batchSize && i+j < totalUsers; j++ {
			username := usernames[i+j]
			city := cities[rand.Intn(len(cities))]
			points := rand.Intn(1_000_001)
			args = append(args, username, city, points)
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("($%d, $%d, $%d)", j*3+1, j*3+2, j*3+3))
		}
		_, err := db.Exec(sb.String(), args...)
		if err != nil {
			log.Fatalf("Failed to insert users: %v", err)
		}
		if (i/batchSize)%100 == 0 {
			log.Printf("Inserted %d users", i+batchSize)
		}
	}
	log.Println("Users generated!")

	// 2. Заполнить friends (user_id уже автонумеруется с 1)
	log.Println("Generating friends...")
	for i := 1; i <= totalUsers; i++ {
		friendCount := rand.Intn(maxFriends-minFriends+1) + minFriends
		friendSet := make(map[int]struct{})
		for len(friendSet) < friendCount {
			f := rand.Intn(totalUsers) + 1 // от 1 до totalUsers
			if f != i {
				friendSet[f] = struct{}{}
			}
		}
		// Вставлять сразу батчем по пользователю
		if len(friendSet) == 0 {
			continue
		}
		var sb strings.Builder
		sb.WriteString("INSERT INTO friends (user_id, friend_id) VALUES ")
		args := make([]interface{}, 0, len(friendSet)*2)
		cnt := 0
		for fid := range friendSet {
			if cnt > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("($%d, $%d)", cnt*2+1, cnt*2+2))
			args = append(args, i, fid)
			cnt++
		}
		_, err := db.Exec(sb.String(), args...)
		if err != nil {
			log.Fatalf("Failed to insert friends for user %d: %v", i, err)
		}
		if i%100_000 == 0 {
			log.Printf("Inserted friends for %d users", i)
		}
	}
	log.Println("Friends generated!")
}
