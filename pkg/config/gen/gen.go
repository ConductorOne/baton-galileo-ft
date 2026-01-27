package main

import (
	cfg "github.com/conductorone/baton-galileo-ft/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("galileoft", cfg.Configuration)
}
