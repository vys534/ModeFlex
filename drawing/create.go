package drawing

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"ModeFlex/utils"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"image"
	"os"
	"time"
)

// TODO: move these somewhere better
var height = 80
var width = 400
var bgBaseOriginX = 80.0
var bgBaseOriginY = 40.0
var iconOffsetX = 3
var iconOffsetY = 3
var globalRankOffsetX = 70
var globalRankOffsetY = 16
var countryRankOffsetY = 22

func CreateCard(userID int64, banner string, userModeData []*api.User) error {

	userImage, err := gg.LoadImage(fmt.Sprintf("./assets/generated/user_images/%d.jpg", userID))
	if err != nil {
		return errors.New(fmt.Sprintf("error while fetching user ID %d avatar from filesystem: %v", userID, err))
	}
	userImageResized := image.NewRGBA(image.Rect(0, 0, 70, 70))
	draw.ApproxBiLinear.Scale(userImageResized, userImageResized.Rect, userImage, userImage.Bounds(), draw.Over, nil)

	bannerExists := false
	// TODO: give users custom banner support in the future
	customBanner, err := gg.LoadImage(fmt.Sprintf("./assets/banners/%s.png", banner))
	if err != nil {
		if os.IsNotExist(err) {
		} else {
			return errors.New(fmt.Sprintf("error while loading banner: %v", err))
		}
	} else {
		bannerExists = true
	}

	c := gg.NewContext(width, height)

	// Draw Banner
	if bannerExists {
		c.DrawRoundedRectangle(0, 0, float64(c.Width()), float64(c.Height()), 5)
		c.Clip()
		c.DrawImageAnchored(customBanner, c.Width()/2, c.Height()/2, 0.5, 0.5)
		c.ResetClip()
		c.SetRGBA(0, 0, 0, 0.7)
	} else {
		c.SetRGBA(0, 0, 0, 1)
	}

	// Draw Black Background
	c.DrawRoundedRectangle(0, 0, float64(c.Width()), float64(c.Height()), 5)
	c.Fill()

	// Draw Clip for avatar
	c.DrawRoundedRectangle(5, 5, 70, 70, 5)
	c.Clip()
	c.DrawImageAnchored(userImageResized, 40, 40, 0.5, 0.5)
	c.ResetClip()

	// Draw country flag
	iconCountry, err := gg.LoadImage(fmt.Sprintf("./assets/flags/%s.png", userModeData[0].CountryCode))
	if err != nil {
		return errors.New(fmt.Sprintf("error while loading flag id %s: %v", userModeData[0].CountryCode, err))
	}
	iconCountryResize := image.NewRGBA(image.Rect(0, 0, 20, 11))
	draw.ApproxBiLinear.Scale(iconCountryResize, iconCountryResize.Rect, iconCountry, iconCountry.Bounds(), draw.Over, nil)
	c.DrawImageAnchored(iconCountryResize, 85, 23, 0, 0)

	// Draw Username
	c.SetHexColor("#ffffff")
	err = c.LoadFontFace("./assets/font/Aldrich-Regular.ttf", 18)
	if err != nil {
		return errors.New(fmt.Sprintf("error while loading font: %v", err))
	}
	c.DrawStringAnchored(userModeData[0].Username, 113, 20, 0, 1)

	// Draw each mode background box
	// TODO: allow users to choose these RGB values
	c.SetRGBA(46.0/255.0, 0, 0, 0.8)
	for i := 0; i < len(userModeData); i++ {
		c.DrawRoundedRectangle(bgBaseOriginX*float64(i+1), bgBaseOriginY, 75, 35, 5)
		c.Fill()
	}

	// Draw ranks
	for i, modeData := range userModeData {
		var icon *image.RGBA

		switch i {
		case 0:
			icon = data.IconOsu
			break
		case 1:
			icon = data.IconTaiko
			break
		case 2:
			icon = data.IconFruits
			break
		case 3:
			icon = data.IconMania
			break
		default:
			icon = data.IconOsu
		}

		c.DrawImageAnchored(icon, int(bgBaseOriginX)*(i+1)+iconOffsetX, int(bgBaseOriginY)+iconOffsetY, 0, 0)

		// Draw Rank Text
		err = c.LoadFontFace("./assets/font/Aldrich-Regular.ttf", 14)
		if err != nil {
			return errors.New(fmt.Sprintf("error while loading font: %v", err))
		}
		// TODO: allow users to choose these RGB values
		c.SetRGBA(255.0/255.0, 153.0/255.0, 153.0/255.0, 1)
		c.DrawStringAnchored(utils.FormatRank(modeData.Statistics.GlobalRank), bgBaseOriginX*float64(i+1)+float64(globalRankOffsetX), bgBaseOriginY+float64(globalRankOffsetY), 1, 0)

		// Draw Country Rank Text
		err = c.LoadFontFace("./assets/font/Aldrich-Regular.ttf", 11)
		if err != nil {
			return errors.New(fmt.Sprintf("error while loading font: %v", err))
		}
		// TODO: allow users to choose these RGB values
		c.SetRGBA(255.0/255.0, 191.0/255.0, 191.0/255.0, 1)
		c.DrawStringAnchored(utils.FormatRank(modeData.Statistics.CountryRank), bgBaseOriginX*float64(i+1)+float64(globalRankOffsetX), bgBaseOriginY+float64(countryRankOffsetY), 1, 1)
	}

	// Draw Timestamp
	err = c.LoadFontFace("./assets/font/Aldrich-Regular.ttf", 10)
	if err != nil {
		return errors.New(fmt.Sprintf("error while loading font: %v", err))
	}
	c.SetHexColor("#ffffff")
	c.DrawStringAnchored(fmt.Sprintf("Generated %s UTC", time.Now().UTC().Format("01/02/2006 15:04")), float64(c.Width())-9, 15, 1, 0)

	f, err := os.OpenFile(fmt.Sprintf("./assets/generated/user_cards/%d.png", userID), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("error while opening file for write: %v", err))
	}
	err = c.EncodePNG(f)
	if err != nil {
		return errors.New(fmt.Sprintf("error while encoding png: %v", err))
	}
	err = f.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("error while closing file: %v", err))
	}

	return nil
}
