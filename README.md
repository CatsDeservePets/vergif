# vergif

`vergif` turns images into animated GIFs.

## Installation

```shell
go install github.com/CatsDeservePets/vergif@latest
```

## Usage

```
usage: vergif [-D] [-d delay] [-l count] [-p palette] -o output image ...
  -D	disable Floyd-Steinberg dithering when quantising true-colour images
  -d delay
    	per-frame delay in 1/100 of a second (default 80)
  -l count
    	animation loop count; 0 means forever, -1 means no looping (default 0)
  -o output
    	output file
  -p palette
    	palette for true-colour quantisation; must be plan9 or websafe (default plan9)
```

## Example

```shell
$ vergif -o my.gif Screenshot_01.png Screenshot_02.png Screenshot_03.png
```
