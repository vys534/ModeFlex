package http_client

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func SaveUserAvatar(userID int64) error {
	req, _ := http.NewRequest(http.MethodGet, api.OsuAvatarEndpoint+strconv.Itoa(int(userID))+"?0.png", nil)
	res, err := data.OsuHTTPClient.Do(req)

	if err != nil {
		return errors.New(fmt.Sprintf("error while fetching user ID %d avatar: %v", userID, err))
	}

	avatarFile, err := os.OpenFile(fmt.Sprintf("./assets/generated/user_images/%d.jpg", userID), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open avatar file for user ID %d: %v", userID, err))
	}

	_, err = io.Copy(avatarFile, res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to write the avatar of user ID %d: %v", userID, err))
	}

	return nil
}
