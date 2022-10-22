package http_client

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetUser(userID int64, mode string) (*api.User, error) {
	req, err := APIRequestWithHeaders(http.MethodGet, api.OsuAPIV2Endpoint+api.OsuAPIV2GetUser+strconv.Itoa(int(userID))+"/"+mode+"?key=id", nil)
	if err != nil {
		return nil, err
	}

	res, err := data.OsuHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var u *api.User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, err
	}

	if len(u.CountryCode) == 0 {
		u.CountryCode = "aq"
	}

	u.CountryCode = strings.ToLower(u.CountryCode)
	return u, nil
}

func GetUserAllModeData(userID int64) ([]*api.User, error) {
	var userDataArray = make([]*api.User, 0)

	for _, mode := range data.ModeArray {
		userData, err := GetUser(userID, mode)
		if err != nil {
			return nil, err
		}
		userDataArray = append(userDataArray, userData)
	}

	return userDataArray, nil
}
