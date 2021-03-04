package helpers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func getPhone() string {

	prefix := []string{"+2557", "+2536", "+2547", "+2119", "+2568", "+2541"}

	net := []string{
		"16", "17", "18", "20", "21", "22", "23", "25", "96", "27",
	}

	rand.Seed(time.Now().UnixNano())
	destCtry := prefix[rand.Intn(len(prefix))]
	destNet := net[rand.Intn(len(net))]
	randNum := strconv.Itoa(111111 + rand.Intn(999999-111111))

	return destCtry + destNet + randNum
}

func GetRecipients(
	count int,
	isMulti bool,
) []map[string]string {
	var requests []map[string]string
	for i := 0; i < count; i++ {
		request := map[string]string{"phone": getPhone()}
		if isMulti {
			request["message"] = fmt.Sprintf("Hello to you %v", i)
		}
		requests = append(requests, request)
	}
	return requests
}
