package main

import "github.com/garyburd/redigo/redis"
import "fmt"

// https://godoc.org/github.com/garyburd/redigo/redis#Pool
func newPool() *redis.Pool {
return &redis.Pool{
            MaxIdle: 100, // Maximum number of idle connections
            MaxActive: 32000, // Maximum number of connections allocated
            Dial: func () (redis.Conn, error) {
                    c, err := redis.Dial("tcp", "127.0.0.1:6379")
                    if err != nil {
                       	fmt.Println(err)
                    }
                    // TODO: add redis AUTH for prod
                    return c, err
            },
    }
}

var pool = newPool()

func RegisterUser(username string, password string) string {
	c := pool.Get()
	defer c.Close()
	exists, err := c.Do("HGET", username, "password")
	if err != nil {
		fmt.Println(err)
	}
	if exists != nil {
		fmt.Println(exists)
		return ""
	}
	_, err = c.Do("HSET", username, "password", password)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Printf("Registered %s. Logging in.\n", username)
	return LoginUser(username, password)
}

func LoginUser(username string, password string) string {
	c := pool.Get()
	defer c.Close()
	// Could have used boolean Reply Helper here, but instead did
	// it this way to minimize redis calls
	exists, err := redis.String(c.Do("HGET", username, "password"))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if password == exists {
		if GetToken(username) == "" {
			fmt.Println("Generate new token")
			SetToken(username)
		}
		return fmt.Sprintf("Bearer %s", GetToken(username))
	} else {
		return ""
	}
}

func GetToken(username string) string {
	c := pool.Get()
	defer c.Close()
	getToken, err := redis.String(c.Do("HGET", username, "token"))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Token does not exist yet")
		return ""
	}
	return getToken
}

func SetToken(username string) {
	c := pool.Get()
	defer c.Close()
	token := NewToken(username)
	_, err := c.Do("HSET", username, "token", token)
	if err != nil {
		fmt.Println(err)
	}
}

func DelToken(username string, token string) {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("HDEL", username, "token")
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: Implement
func UpdateEventually(packageName string, deviceID string) {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("SADD", packageName, deviceID)
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: Impelement
func CheckUpdateNeeded(packageName string, deviceID string) {
	c := pool.Get()
	defer c.Close()
	array, err := redis.Strings(c.Do("SMEMBERS", packageName, 0, -1))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(array)
}

func SetInfo(packageHash string, packageName string, fileName string) {
	c := pool.Get()
	defer c.Close()
	// Do combines Send, Flush, and Receive. We need to do manual pipelining
	// for this, we use "Send", then "Flush", we don't need "Receive"
	err := c.Send("HSET", packageHash, "package", packageName)
	if err != nil {
		fmt.Println("Error setting package", err)
	}
	err = c.Send("HSET", packageHash, "link", fileName)
	if err != nil {
		fmt.Println("Error setting hash", err)
	}
	err = c.Flush()
	if err != nil {
		fmt.Println("Error flushing stream", err)
	}
}

func GetInfo(packageHash string) []string {
	fmt.Println(packageHash)
	c := pool.Get()
	defer c.Close()
	array, err := redis.Strings(c.Do("HGETALL", packageHash))
	if len(array) >= 4 {
		packageName, err := redis.String(array[1], nil)
		hash, err := redis.String(array[3], nil)
		if err != nil {
			return nil
		}
		downloadLink := BASE_URL + "/apks/" + hash
		return []string{packageName, downloadLink}
	}
	if err != nil {
		fmt.Println(err)
	}
	return nil
}