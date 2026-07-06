package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type paletteKind byte

const (
	unknown paletteKind = iota
	plan9
	webSafe
	// TODO: auto
)

func (p paletteKind) String() string {
	switch p {
	case plan9:
		return "plan9"
	case webSafe:
		return "websafe"
	default:
		return "unknown"
	}
}

func (p *paletteKind) Set(s string) error {
	switch strings.ToLower(s) {
	case "plan9":
		*p = plan9
	case "websafe":
		*p = webSafe
	default:
		return errors.New("must be plan9 or websafe")
	}
	return nil
}

func (p paletteKind) Palette() color.Palette {
	switch p {
	case plan9:
		return palette.Plan9
	case webSafe:
		return palette.WebSafe
	default:
		panic(fmt.Sprintf("unknown palette: %d", p))
	}
}

var (
	noDither  = flag.Bool("D", false, "disable Floyd-Steinberg dithering when quantising true-colour images")
	delay     = flag.Uint("d", 80, "per-frame `delay` in 1/100 of a second")
	loopCount = flag.Int("l", 0, "animation loop `count`; 0 means forever, -1 means no looping (default 0)")
	outPath   = flag.String("o", "", "`output` file")
	pal       = plan9
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: vergif [-D] [-d delay] [-l count] [-p palette] -o output image ...")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("vergif: ")
	flag.Usage = usage
	flag.Var(&pal, "p", "`palette` for true-colour quantisation; must be plan9 or websafe")
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
			log.Fatalf("%s: %v", path, err)
		}

		anim.Image = append(anim.Image, palettise(img, pal.Palette(), !*noDither))
		anim.Disposal = append(anim.Disposal, gif.DisposalBackground)
		anim.Delay = append(anim.Delay, int(*delay))
	}

	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := gif.EncodeAll(f, anim); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func palettise(img image.Image, pal color.Palette, dither bool) *image.Paletted {
	if p, ok := img.(*image.Paletted); ok {
		return p
	}

	var drawer draw.Drawer = draw.Src
	if dither {
		drawer = draw.FloydSteinberg
	}

	frame := image.NewPaletted(img.Bounds(), pal)
	drawer.Draw(frame, frame.Rect, img, img.Bounds().Min)
	return frame
}
