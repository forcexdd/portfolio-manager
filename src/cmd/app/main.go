package main

import "os"

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "web":
		runWeb()
	case "desktop":
		runDesktop()
	default:
		os.Exit(1)
	}
}
