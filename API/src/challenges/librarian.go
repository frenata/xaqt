package challenges

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var challenges map[Id]Challenge

func Get() Challenge {
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(len(challenges))
	id := strconv.Itoa(n)
	return challenges[id]

	// return challenges["1"]
}

func GetById(i Id) Challenge {
	return challenges[i]
}

func GetByTag() {

}

func GetAll() map[Id]Challenge {
	return challenges
}

func init() {
	// This should connect to SQL database and perform
	// some check to verify connection

	log.Println("Reading challenges file...")
	bytes, err := ioutil.ReadFile("data/challenges.json")
	if err != nil {
		panic(err)
	}

	challenges = make(map[Id]Challenge)
	err = json.Unmarshal(bytes, &challenges)
	if err != nil {
		panic(err)
	}
	log.Println("Challenges file loaded.")
	// for k, v := range challenges {
	// 	log.Printf("Id: %s maps to %s", k, v.Id)
	// }
}
