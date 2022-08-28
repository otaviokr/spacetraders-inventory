package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	httpEndpointGetUserDetails = "https://api.spacetraders.io/my/account?token=%s"
	httpEndpointGetShipList    = "https://api.spacetraders.io/my/ships?token=%s"
	httpEndpointGetLeaderboard = "https://api.spacetraders.io/game/leaderboard/net-worth?token=%s"
	httpEndpointGetGamestatus  = "https://api.spacetraders.io/game/status?token=%s"

	maxRetriesTimeout = 5
	waitTimeout       = time.Duration(10)
)

// WebProxy is an implementation of web.Proxy.
type WebProxy struct {
	token   string
	baseUrl string
}

// NewWebProxy creates a new instance of WebProxy.
//
// token is provided by the game when you claim your username.
func NewWebProxy(token string) Proxy {
	return &WebProxy{
		token:   token,
		baseUrl: ""}
}

// GetUserDetails collects user details.
//
// https://api.spacetraders.io/#api-account-GetAccount
func (wp *WebProxy) GetUserDetails() ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetUserDetails, wp.token))
}

// GetShipList collects the list of ships the user owns.
//
// https://api.spacetraders.io/#api-ships-GetShips
func (wp *WebProxy) GetShipList() ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetShipList, wp.token))
}

// GetLeaderboard collects the leaderboard and the user rank.
//
// https://api.spacetraders.io/#api-leaderboard-netWorth
func (wp *WebProxy) GetLeaderboard() ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetLeaderboard, wp.token))
}

// GetGameStatus collects the game status (if it is online and accessible).
//
// https://api.spacetraders.io/#api-game-Status
func (wp *WebProxy) GetGameStatus() ([]byte, error) {
	return wp.get(fmt.Sprintf(httpEndpointGetGamestatus, wp.token))
}

// get is a generic GET request, used by the other applications.
func (wp *WebProxy) get(uri string) ([]byte, error) {
	count := 0
	for count < maxRetriesTimeout {
		response, err := http.Get(uri)
		if err == nil {
			defer response.Body.Close()
			return io.ReadAll(response.Body)
		}

		errUrl := err.(*url.Error)
		if errUrl.Timeout() {
			time.Sleep(waitTimeout * time.Second)
		} else {
			return []byte{}, err
		}
	}
	return []byte{}, fmt.Errorf("reached unexpected piece of code in get")
}

// // post is a generic POST request, used by the other applications.
// func (wp *WebProxy) post(uri string, data io.Reader) ([]byte, error) {
// 	count := 0
// 	for count < MAX_RETRIES_TIMEOUT {
// 		response, err := http.Post(uri, "application/json", data)
// 		if err == nil {
// 			defer response.Body.Close()
// 			return io.ReadAll(response.Body)
// 		}

// 		errUrl := err.(*url.Error)
// 		if errUrl.Timeout() {
// 			time.Sleep(WAIT_TIMEOUT * time.Second)
// 		} else {
// 			return []byte{}, err
// 		}
// 	}
// 	return []byte{}, fmt.Errorf("reached unexpected piece of code in post")
// }
