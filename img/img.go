package img

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ImageControl struct {
	name string      // 文件路径
	img  image.Image // 原始图片
}

func (ic *ImageControl) LoadImageForImageObj(img *image.Image) (err error) {
	ic.img = *img
	return
}

func (ic *ImageControl) LoadImageForFile(name string) (err error) {
	ic.name = name
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()
	ic.img, _, err = image.Decode(file)
	return
}

func (ic *ImageControl) CopyImage(position [][]float64, offset []float64) (err error) {
	ptsF := []float64{
		position[0][0] + offset[0],
		position[0][1] + offset[1],
		position[2][0] + offset[2],
		position[2][1] + offset[3],
	}
	pts := []int{int(ptsF[0]), int(ptsF[1]), int(ptsF[2]), int(ptsF[3])}
	log.Printf("TrimmingForFile: pts: %#v => %#v, ", ptsF, pts)
	if rgbImg, ok := ic.img.(*image.YCbCr); ok {
		ic.img = rgbImg.SubImage(image.Rect(pts[0], pts[1], pts[2], pts[3])).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := ic.img.(*image.RGBA); ok {
		ic.img = rgbImg.SubImage(image.Rect(pts[0], pts[1], pts[2], pts[3])).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := ic.img.(*image.NRGBA); ok {
		ic.img = rgbImg.SubImage(image.Rect(pts[0], pts[1], pts[2], pts[3])).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {
		err = errors.New("图片解码失败")
	}
	return
}

func (ic *ImageControl) SaveImage(destFile string) (err error) {
	f, err := os.OpenFile(destFile, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(destFile)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		err = jpeg.Encode(f, ic.img, &jpeg.Options{Quality: 80})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, ic.img)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, ic.img, &gif.Options{NumColors: 256})
	}
	return err
}

func TrimmingForFile(srcFile string, destFile string, position [][]float64, offset []float64) (err error) {
	var ic ImageControl
	log.Printf("TrimmingForFile: file: %s => %s, \n Pos: %#v\n off: %#v", srcFile, destFile, position, offset)
	err = ic.LoadImageForFile(srcFile)
	if err != nil {
		return err
	}
	err = ic.CopyImage(position, offset)
	if err != nil {
		return err
	}
	err = ic.SaveImage(destFile)
	return
}

func TrimmingForImageObj(img *image.Image, destFile string, position [][]float64, offset []float64) (err error) {
	var ic ImageControl
	err = ic.LoadImageForImageObj(img)
	if err != nil {
		return err
	}
	err = ic.CopyImage(position, offset)
	if err != nil {
		return err
	}
	err = ic.SaveImage(destFile)
	return
}
