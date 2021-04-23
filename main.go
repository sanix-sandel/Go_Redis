package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	album := new(Album)
	album.Title = reply["title"]
	album.Artist = reply["artist"]

	album.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}

	album.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil, err
	}
	return album, nil
}

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
	}
	//get titile and converting it to string
	title, err := redis.String(conn.Do("HGET", "album:2", "title"))
	if err != nil {
		log.Fatal(err)
	}*/
	//HGETALL to get array of replies
	reply, err := redis.StringMap(conn.Do("HGETALL", "album:2"))
	if err != nil {
		log.Fatal(err)
	}
	album, err := populateAlbum(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(album)

	//fmt.Println("Electric Ladyland added")

}
