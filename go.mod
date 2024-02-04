module github.com/filevich/truco-ai

go 1.18

require github.com/truquito/truco v0.1.0

require (
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

require (
	github.com/filevich/canvas v0.0.0 // indirect
	github.com/filevich/combinatronics v0.0.0-20220316214652-26aa6db09482
	github.com/jedib0t/go-pretty/v6 v6.5.4
)

replace github.com/truquito/truco => ../truco
