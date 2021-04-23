package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

func main() {
	//establish connection
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	//close the connection
	defer conn.Close()

	/*_, err = conn.Do("HMSET", "album:2", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	if err != nil {
		log.Fatal(err)
	}*/
	//get titile and converting it to string
	title, err := redis.String(conn.Do("HGET", "album:2", "title"))
	if err != nil {
		log.Fatal(err)
	}
	artist, err := redis.String(conn.Do("HGET", "album:2", "artist"))
	if err != nil {
		log.Fatal(err)
	}
	price, err := redis.Float64(conn.Do("HGET", "album:2", "price"))
	if err != nil {
		log.Fatal(err)
	}
	likes, err := redis.Int(conn.Do("HGET", "album:2", "likes"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s by %s: $%.2f [%d likes]\n", title, artist, price, likes)

	//fmt.Println("Electric Ladyland added")

}
