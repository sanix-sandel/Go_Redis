package main

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

var ErrNoAlbum = errors.New("no album found")

type Album struct {
	Title  string  `redis:"title"`
	Artist string  `redis:"artist"`
	Price  float64 `redis:"price"`
	Likes  int     `redis:"likes"`
}

func FindAlbum(id string) (*Album, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn := pool.Get()

	// Importantly, use defer and the connection's Close() method to
	// ensure that the connection is always returned to the pool before
	// FindAlbum() exits.
	defer conn.Close()

	// Fetch the details of a specific album. If no album is found
	// the given id, the []interface{} slice returned by redis.Values
	// will have a length of zero. So check for this and return an
	// ErrNoAlbum error as necessary.
	values, err := redis.Values(conn.Do("HGETALL", "album:"+id))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, ErrNoAlbum
	}

	var album Album
	err = redis.ScanStruct(values, &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func IncrementLikes(id string) error {
	conn := pool.Get()
	defer conn.Close()

	//check that an album with the given id exists
	exists, err := redis.Int(conn.Do("EXISTS", "album:"+id))
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	//Using Multi to inform Redis that we are starting a new transaction
	//The conn.Send() method writes the command to the
	// connection's output buffer -- it doesn't actually send it to the
	// Redis server... despite it's name!
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}

	//Icrement the number of likes in the album hash by 1.
	// Because follows a MULTI command, this HINCRBY command is NOT executed but
	// it is QUEUED as part of the transaction. We still need to check
	// the reply's Err field at this point in case there was a problem
	// queueing the command.

	err = conn.Send("HINCRBY", "album:"+id, "likes", 1)
	if err != nil {
		return err
	}
	// And we do the same with the increment on our sorted set.
	err = conn.Send("ZINCRBY", "likes", 1, id)
	if err != nil {
		return err
	}

	// Execute both commands in our transaction together as an atomic
	// group. EXEC returns the replies from both commands but, because
	// we're not interested in either reply in this example, it
	// suffices to simply check for any errors. Note that calling the
	// conn.Do() method flushes the previous commands from the
	// connection output buffer and sends them to the Redis server.
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}

	return nil

}
