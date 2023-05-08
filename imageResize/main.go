package main

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

func main() {

	args := os.Args
	var desired_width, desired_height, scala int
	var err error
	if len(args) == 4 {
		desired_width, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
			return
		}
		desired_height, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
			return
		}
		scala, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		desired_width = -1
		desired_height = -1
		scala = 85
	}

	imageList := []string{}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpg") && !strings.HasPrefix(file.Name(), "copy_") {
			imageList = append(imageList, file.Name())
		}
	}

	for _, imageName := range imageList {

		imageResize(imageName, "copy_"+imageName, desired_width, desired_height, scala)
	}
}

func imageResize(inputImagePath, outputImagePath string, desiredWidth, desiredHeight, scala int) {
	inputFile, err := os.Open(inputImagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// 解码图像
	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// 获取原始图像的大小
	inputBounds := inputImage.Bounds()
	width := inputBounds.Max.X
	height := inputBounds.Max.Y

	var newWidth, newHeight int
	if desiredHeight != -1 {
		if width <= desiredWidth && height <= desiredHeight {
			log.Println("Image is already smaller than desired size")
			return
		}

		// 根据指定的大小计算新的图像大小

		if width > height {
			newWidth = desiredWidth
			newHeight = int(float64(height) * (float64(desiredWidth) / float64(width)))
		} else {
			newHeight = desiredHeight
			newWidth = int(float64(width) * (float64(desiredHeight) / float64(height)))
		}
	} else {
		newWidth = width
		newHeight = height
	}
	// 如果原始图像已经小于或等于指定的大小，则无需进行压缩

	// 缩放图像
	outputImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(outputImage, outputImage.Bounds(), inputImage, inputBounds, draw.Over, nil)

	// 创建输出文件
	//outputImagePath := "output_image.jpg"
	outputFile, err := os.Create(outputImagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 将图像编码为JPEG格式并保存
	options := &jpeg.Options{Quality: scala}
	err = jpeg.Encode(outputFile, outputImage, options)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Image saved to", outputImagePath)
}
