package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"golang.org/x/image/draw"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	var images []image.Image
	for _, arg := range os.Args[1:] {
		// Skip screenshot flag
		if arg == "-screenshot" {
			continue
		}
		var img image.Image
		var err error

		// Try HTTP first
		resp, httpErr := http.Get(arg)
		if httpErr == nil {
			defer resp.Body.Close()
			img, _, err = image.Decode(resp.Body)
		} else {
			// Fallback to local file
			f, openErr := os.Open(arg)
			if openErr != nil {
				log.Printf("failed to fetch/open image %s: %v / %v", arg, httpErr, openErr)
				continue
			}
			defer f.Close()
			img, _, err = image.Decode(f)
		}

		if err != nil {
			log.Printf("failed to decode image %s: %v", arg, err)
			continue
		}
		images = append(images, img)
	}
	if len(images) == 0 {
		// Load the local GOTUI logo - try multiple paths
		logoPaths := []string{"logo.png", "../../logo.png", "_examples/image/../../logo.png"}
		var img image.Image
		var loaded bool
		for _, path := range logoPaths {
			f, err := os.Open(path)
			if err != nil {
				continue
			}
			img, _, err = image.Decode(f)
			f.Close()
			if err == nil {
				loaded = true
				break
			}
		}
		if !loaded {
			log.Fatalf("failed to open logo.png from any path")
		}
		// Resize the logo to fit the widget inner dimensions (widget - borders)
		img = resizeImage(img, 118, 33)
		images = append(images, img)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	img := widgets.NewImage(nil)
	img.SetRect(0, 0, 120, 35)
	index := 0
	render := func() {
		img.Image = images[index]
		if !img.Monochrome {
			img.Title = fmt.Sprintf("Color %d/%d", index+1, len(images))
		} else if !img.MonochromeInvert {
			img.Title = fmt.Sprintf("Monochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		} else {
			img.Title = fmt.Sprintf("InverseMonochrome(%d) %d/%d", img.MonochromeThreshold, index+1, len(images))
		}
		ui.Render(img)
	}
	render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if handleImageEvents(img, e, &index, len(images)) {
			return
		}
		render()
	}
}

// resizeImage scales an image to the target width and height using high-quality interpolation
func resizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func handleImageEvents(img *widgets.Image, e ui.Event, index *int, count int) bool {
	switch e.ID {
	case "q", "<C-c>":
		return true
	case "<Left>", "h":
		*index = (*index + count - 1) % count
	case "<Right>", "l":
		*index = (*index + 1) % count
	case "<Up>", "k":
		img.MonochromeThreshold++
	case "<Down>", "j":
		img.MonochromeThreshold--
	case "<Enter>":
		img.Monochrome = !img.Monochrome
	case "<Tab>":
		img.MonochromeInvert = !img.MonochromeInvert
	}
	return false
}
