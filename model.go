package main

import (
	"errors"
	"github.com/spf13/cast"
	"image"
	"strings"
)

type PlistInfo struct {
	Frames   map[string]PlistFrame `plist:"frames"`
	Metadata PlistMeta             `plist:"metadata"`
}

type PlistFrame struct {
	Frame           Rectangle `plist:"frame"`
	Offset          Point     `plist:"offset"`
	Rotated         bool      `plist:"rotated"`
	SourceColorRect Rectangle `plist:"sourceColorRect"`
	SourceSize      Point     `plist:"sourceSize"`
}

func (f *PlistFrame) GetSourceRect() image.Rectangle {
	if f.Rotated {
		r := f.SourceColorRect.Rectangle
		return image.Rect(r.Min.X, r.Min.Y, r.Max.Y, r.Max.X)
	} else {
		return f.SourceColorRect.Rectangle
	}
}
func (f *PlistFrame) GetSourceSize() image.Rectangle {
	return image.Rect(0, 0, f.SourceSize.X, f.SourceSize.Y)
}

type PlistMeta struct {
	Format              int    `plist:"format"`
	RealTextureFileName string `plist:"realTextureFileName"`
	Size                string `plist:"size"`
	TextureFileName     string `plist:"textureFileName"`
}

type Point struct {
	image.Point
}

func (p *Point) UnmarshalPlist(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	items := strings.Split(str, ",")
	if len(items) != 2 {
		return errors.New("not point string")
	}
	p.X = cast.ToInt(items[0])
	p.Y = cast.ToInt(items[1])
	return nil
}

type Rectangle struct {
	image.Rectangle
}

func (p *Rectangle) UnmarshalPlist(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	items := strings.Split(str, ",")
	if len(items) != 4 {
		return errors.New("not rect string")
	}
	p.Min.X = cast.ToInt(items[0])
	p.Min.Y = cast.ToInt(items[1])
	p.Max.X = cast.ToInt(items[2])
	p.Max.Y = cast.ToInt(items[3])
	return nil
}
