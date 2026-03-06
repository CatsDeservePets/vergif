package main

import (
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var (
	progName  = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	delay     = flag.Uint("delay", 80, "delay per frame in 1/100 of a second")
	loopCount = flag.Int("loop", 0, "animation loop count; 0 means forever, -1 means no looping (default 0)")
	outPath   = flag.String("o", "", "`output` file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s -o output [flags] image ...\n", progName)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 || *outPath == "" {
		flag.Usage()
	}

	anim := &gif.GIF{
		Image:     make([]*image.Paletted, 0, flag.NArg()),
		Delay:     make([]int, 0, flag.NArg()),
		LoopCount: *loopCount,
		Disposal:  make([]byte, 0, flag.NArg()),
	}

	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			log.Fatal(err)
		}

		var frame *image.Paletted
		if p, ok := img.(*image.Paletted); ok {
			frame = p
		} else {
			frame = image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.FloydSteinberg.Draw(frame, frame.Rect, img, img.Bounds().Min)
		}

		anim.Image = append(anim.Image, frame)
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)
		anim.Delay = append(anim.Delay, int(*delay))
	}

	f, err := os.OpenFile(*outPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := gif.EncodeAll(f, anim); err != nil {
		log.Fatal(err)
	}
}
