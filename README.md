# pymosaic
Python wrapper around go code to generate photomosaics

An interesting feature in Go 1.5 is build-modes flag during the build phase. It's possible to compile shared library
ans use that from your python code.

This is a thin wrapper aroung go library for photomosaic generation.

## Installation
From the src directory

pymosaic/src$ go build -buildmode=c-shared -o ../dist/pymosaic.so 

Now copy pymosaic.so to your virtual environments site-packages.

## Quick Start
You can read about photomosaics here:
https://en.wikipedia.org/wiki/Photographic_mosaic
You need a library of tile images and a target photo to generate a photomosaic.

### flickr client:
pymosaic also includes a flickr client to download images with sepcified tag to help build your tile library.

##### pymosaic.download
```python
download(flickr_tag, number, libdir, api_key)
```
flickr_tag(string): the images tagged with this term will be downloaded
number(int): the number of images that are to be downloaded
libdir(string): the output directory where the images should be downloaded
api_key(string): your flickr api key

example: 
```python
>>>import pymosaic
>>>pymosaic.download("waffles", 200, "<homedir>/tiledir", "<flickr api key>")
```
### mosaic generation

##### pymosaic.analyze_images
```python
analyze_images(libdir)
```
libdir(string): the path to tile image library
this function takes path of the tile directory as input and returns a json string of image statistics
(to be used in next step of mosaic generation

example:
```python
>>>img_statistics = pymosaic.analyze_images("<homedir>/tiledir")
```

##### pymosaic.generate_mosaic

```python
generate_mosaic(<target_image>, <mosaic_image>, <img_statistics>)
```
target_image(string): The path of image that you want to create mosaic out of.
mosaic_image(string): the path of the mosaic image to be generated (should be .png and you need to have write permissions)
img_statistics(string): image statistics json string generate in previous step (analyze_images)

example:
```python
>>>pymosaic.generate_mosaic("/home/user/my_profile.jpg",
                           "/home/user/output.png", img_statistics)
```
