package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewImageOptions(t *testing.T) {
	img, err := NewImageOptions("test_image.jpg")
	assert.Nil(t, err)

	expected := ImageOptions{
		UserId:1, //temp
		FileName: "test_image",
		Format: "jpg",
	}

	assert.Equal(t, expected, img)

	//

	_, err = NewImageOptions("test_file.mp4")
	assert.Equal(t, fmt.Errorf("unsupported image format"), err)
}

func TestImageOptions_Save(t *testing.T) {

	err := initApp()
	assert.Nil(t, err)

	imgRaw, err := os.OpenFile("testdata/test_image.jpg", os.O_RDONLY, 0644)
	assert.Nil(t, err)

	imgOpts := ImageOptions{
		UserId:1, //temp
		FileName: "test_image",
		Format: "jpg",
	}

	_, err = imgOpts.Save(imgRaw)
	assert.Nil(t, err)
}

func TestImageOptions_Resize(t *testing.T) {
	imgRaw, err := os.OpenFile("testdata/test_image.jpg", os.O_RDONLY, 0644)
	assert.Nil(t, err)

	imgOpts := ImageOptions{
		Id:10,
		UserId:1, //temp
		FileName: "test_image",
		Format: "jpg",
	}

	imgResOpts, imgRes, err := imgOpts.Resize(imgRaw, 50, 40)
	assert.Nil(t, err)

	expectedRes := ImageResizedOptions{
		UserId: 1,  //temp
		OriginalId: 10,
		FileName: "test_image",
		Format: "jpg",
		Width:50,
		Height:40,
	}

	assert.Equal(t, expectedRes, imgResOpts)

	///

	b1, err := ioutil.ReadAll(imgRes)
	assert.Nil(t, err)

	f2, err := os.OpenFile("testdata/test_resized.jpg", os.O_RDONLY, 0644)
	assert.Nil(t, err)
	b2, err := ioutil.ReadAll(f2)
	assert.Nil(t, err)

	assert.Equal(t, b1, b2)

}

func TestImageOptions_Raw(t *testing.T) {

	err := initApp()
	assert.Nil(t, err)

	imgRaw, err := os.OpenFile("testdata/test_image.jpg", os.O_RDONLY, 0644)
	assert.Nil(t, err)

	imgOpts := ImageOptions{
		Id:1,
		UserId:1, //temp
		FileName: "test_image",
		Format: "jpg",
	}

	_, err = imgOpts.Save(imgRaw)
	assert.Nil(t, err)

	imgRaw2, err := imgOpts.Raw()
	assert.Nil(t, err)

	b1, err := ioutil.ReadAll(imgRaw)
	assert.Nil(t, err)
	b2, err := ioutil.ReadAll(imgRaw2)
	assert.Nil(t, err)

	assert.Equal(t, b1, b2)

}