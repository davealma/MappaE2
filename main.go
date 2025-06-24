package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Star struct {
	Resonance int `json:"resonance`
}

var galaxies int
func GetOracleInfo(stelarJump int, ignoreCount bool) []Star {
	baseUrl := os.Getenv("API_URL")
	req, err := http.NewRequest("Get", baseUrl + `/v1/s1/e2/resources/stars?page=`+strconv.Itoa(stelarJump)+`&sort-by=id&sort-direction=desc`, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("API-KEY", os.Getenv("API_KEY"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	totalCountHeader := resp.Header.Get("x-total-count");

	var stars []Star

	json.Unmarshal([]byte(body), &stars)
	
	
	if !ignoreCount {
		count, err := strconv.Atoi(totalCountHeader); if err != nil{
			panic(err)
		}	
		galaxies = int(count/3) + 1		
	}

	return stars
}

func PostSolution(avg int) {
	baseUrl := os.Getenv("API_URL")

	payload := []byte(`{
		"average_resonance": "`+strconv.Itoa(avg)+`"
	}`)

	req, err := http.NewRequest("POST", baseUrl+"/v1/s1/e2/solution", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("API-KEY", os.Getenv("API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err!= nil {
		panic(err)
	}
	fmt.Println("Response Status: ", string(bodyResp))
}

func main() {
	godotenv.Load()
	fmt.Println("El cosmos vibra en una sinfonía matemática. La resonancia de cada estrella se construye sobre la anterior, pero el Oráculo te presenta las estrellas en un orden cósmico propio.")
	GetOracleInfo(1, false)
	var all []Star
	for i := 0; i < galaxies; i++ {		
		stars := GetOracleInfo(i+1, true)
		all = append(all, stars...)		
	}
	fmt.Println(len(all))
	
	total := 0
	for i := 0; i < len(all); i++ {
		total = all[i].Resonance + total
	}

	avg := total/len(all)
	fmt.Println("Avg: ", avg)
	PostSolution(avg)
}