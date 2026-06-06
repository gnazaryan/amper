package templates

import "embed"

// content holds our static web server content.
//
//go:embed notification/*
var templates embed.FS

func GetFS() embed.FS {
	return templates
}
