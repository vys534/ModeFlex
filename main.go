package main

import (
	"ModeFlex/data"
	"ModeFlex/http_client"
	"ModeFlex/routing"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"golang.org/x/image/draw"
	"image"
	"log"
	"os"
	"os/signal"
)

var iconSize = 15

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("No config file set, %v", err)
		} else {
			log.Fatalf("Error reading config file, %v", err)
		}
	}

	viper.SetDefault("Port", 3000)
	viper.SetDefault("UpdateInterval", 6000)

	err := viper.Unmarshal(&data.Configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if len(data.Configuration.ClientSecret) == 0 || len(data.Configuration.ClientID) == 0 {
		log.Fatalf("Client ID or client secret are empty!")
	}

	resizeAllStaticImages()

	err = http_client.ClientCredentialsGrant()
	if err != nil {
		log.Fatalf("Client credentials grant errored! %+v", err)
	}

	s := &fasthttp.Server{
		Handler:                       routing.Route,
		Concurrency:                   512,
		GetOnly:                       true,
		DisableKeepalive:              false,
		TCPKeepalive:                  false,
		TCPKeepalivePeriod:            0,
		LogAllErrors:                  false,
		DisableHeaderNamesNormalizing: false,
		NoDefaultServerHeader:         true,
		NoDefaultDate:                 true,
		KeepHijackedConns:             false,
	}

	fmt.Printf("Listening for requests on port %d\n", data.Configuration.Port)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := s.ListenAndServe(fmt.Sprintf(":%d", data.Configuration.Port)); err != nil {
			log.Fatalf("Listen error: %v\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down")

	if err := s.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown gracefully: %v\n", err)
	}

	log.Println("Shut down")
	os.Exit(0)
}

func resizeAllStaticImages() {
	iconOsu, err := gg.LoadImage("./assets/mode_icons/osu.png")
	if err != nil {
		log.Fatal(err)
	}
	iconTaiko, err := gg.LoadImage("./assets/mode_icons/taiko.png")
	if err != nil {
		log.Fatal(err)
	}
	iconFruits, err := gg.LoadImage("./assets/mode_icons/fruits.png")
	if err != nil {
		log.Fatal(err)
	}
	iconMania, err := gg.LoadImage("./assets/mode_icons/mania.png")
	if err != nil {
		log.Fatal(err)
	}

	iconOsuResize := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	iconTaikoResize := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	iconFruitsResize := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	iconManiaResize := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))

	draw.ApproxBiLinear.Scale(iconOsuResize, iconOsuResize.Rect, iconOsu, iconOsu.Bounds(), draw.Over, nil)
	draw.ApproxBiLinear.Scale(iconTaikoResize, iconTaikoResize.Rect, iconTaiko, iconTaiko.Bounds(), draw.Over, nil)
	draw.ApproxBiLinear.Scale(iconFruitsResize, iconFruitsResize.Rect, iconFruits, iconFruits.Bounds(), draw.Over, nil)
	draw.ApproxBiLinear.Scale(iconManiaResize, iconManiaResize.Rect, iconMania, iconMania.Bounds(), draw.Over, nil)

	data.IconOsu = iconOsuResize
	data.IconTaiko = iconTaikoResize
	data.IconFruits = iconFruitsResize
	data.IconMania = iconManiaResize
}
