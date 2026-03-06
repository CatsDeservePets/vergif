# vergif

`vergif` turns images into animated GIFs.

## Installation

```shell
go install github.com/CatsDeservePets/vergif@latest
```

## Usage

```
usage: vergif [flags] -o output image ...
  -delay uint
        delay per frame in 1/100 of a second (default 80)
  -loop int
        animation loop count; 0 means forever, -1 means no looping (default 0)
  -o output
        output file
```

## Example

```shell
$ vergif -o my.gif Screenshot_01.png Screenshot_02.png Screenshot_03.png
```
