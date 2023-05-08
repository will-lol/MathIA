# MathIA
This program was developed for my Math A&A IA to analyse how the filesize of images compressed with JPEG to constant SSIM quality varies as you change the resolution.
## Installation
No dependencies aside from golang.
Clone the repo and run:
```
go build
````
## Usage
Execute the binary and pass in a png image file to compress:
```
./MathIA images/myImage.png
```
Designed for use with GNU parellel. With a folder of png images:
```
parallel -v ./MathIA ::: images/*
```