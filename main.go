package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"howett.net/plist"
	"image"
	"image/draw"
	_ "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please input plist_file or plist_dir")
		os.Exit(-1)
	}
	path := os.Args[1]
	pathStat, err := os.Stat(path)
	if err != nil {
		if pathStat, err = os.Stat(path + ".plist"); err != nil {
			if pathStat, err = os.Stat(path + ".png"); err != nil {
				panic(err)
			}
		}
	}
	if pathStat.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		filenameMap := map[string]os.DirEntry{}
		for _, file := range files {
			items := strings.Split(file.Name(), ".")
			if _, ok := filenameMap[items[0]]; !ok {
				filenameMap[items[0]] = file
			}
		}
		for key, _ := range filenameMap {
			newPath := filepath.Join(path, key)
			if err = handlePath(newPath); err != nil {
				fmt.Printf("handle %s fail result:%s\n", newPath, err)
			} else {
				fmt.Printf("handle %s success\n", newPath)
			}
		}
	} else {
		if err = handlePath(path); err != nil {
			fmt.Printf("handle %s fail result:%s\n", path, err)
		} else {
			fmt.Printf("handle %s success\n", path)
		}
	}
}

func handlePath(path string) error {
	var plistPath string
	var imgPath string
	var outputPath string

	// 判断输入路径合法，以及兼容
	if _, err := os.Stat(path); err == nil || os.IsExist(err) { // 存在
		switch filepath.Ext(path) {
		case ".plist":
			plistPath = path
			outputPath = path[:len(path)-6]
			imgPath = outputPath + ".png"
		case ".png":
			imgPath = path
			outputPath = plistPath[:len(plistPath)-4]
			plistPath = outputPath + ".plist"
		}
	} else { // 不存在
		outputPath = path
		plistPath = path + ".plist"
		imgPath = path + ".png"
	}

	// 判定下两个文件是不是都有
	if _, err := os.Stat(outputPath); err == nil || os.IsExist(err) { // 存在
		if rmErr := os.RemoveAll(outputPath); rmErr != nil {
			return rmErr
		}
		//return errors.New("outputPath is exist:" + outputPath)
	}
	if makeErr := os.MkdirAll(outputPath, fs.ModeDir|fs.ModePerm); makeErr != nil {
		return makeErr
	}
	if _, err := os.Stat(plistPath); !(err == nil || os.IsExist(err)) { // 存在
		return errors.New("no exist:" + plistPath)
	}
	if _, err := os.Stat(imgPath); !(err == nil || os.IsExist(err)) { // 存在
		return errors.New("no exist:" + imgPath)
	}

	// 读取大图片
	f, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	bigImage, formatname, err := image.Decode(f)
	if err != nil {
		fmt.Println(formatname)
		return err
	}

	// 读取文件
	info := PlistInfo{}
	file, err := os.Open(plistPath)
	if err != nil {
		return err
	}
	err = plist.NewDecoder(file).Decode(&info)
	if err != nil {
		return err
	}

	for name, frame := range info.Frames {
		// 先生成个中间的图片
		tmpImg := image.NewNRGBA(frame.GetOriginSize())
		draw.Draw(tmpImg, tmpImg.Bounds(), bigImage, frame.Frame.Min, draw.Over)
		if frame.Rotated {
			tmpImg = imaging.Rotate90(tmpImg)
		}
		// 判断下是否旋转，再贴到目标图片里面
		newImg := image.NewNRGBA(frame.GetSourceSize())
		draw.Draw(newImg, frame.GetTargetRect(), tmpImg, tmpImg.Bounds().Min, draw.Src)
		saveImage(newImg, filepath.Join(outputPath, name))
	}
	return nil
}

func saveImage(img image.Image, path string) {
	outFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("create %s fail result:%s\n", path, err)
		return
	}
	if err = png.Encode(outFile, img); err != nil {
		fmt.Printf("png encode %s fail result:%s\n", path, err)
		return
	}
	if err = outFile.Close(); err != nil {
		fmt.Printf("close %s fail result:%s\n", path, err)
		return
	}
}
