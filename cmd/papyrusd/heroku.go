package main

import "os"

func loadHerokuEnv() {
	if os.Getenv("DYNO") != "" {
		os.Setenv("ADDR", ":"+os.Getenv("PORT"))
		os.Setenv("MONGO_URL", os.Getenv("MONGOLAB_URI"))
	}
}
