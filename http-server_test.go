package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestGetImagesHandler(t *testing.T) {
	err := initApp()
	assert.Nil(t, err)
	router := SetupRouter()

	w := performRequest(router, "GET", "/api/v1/images", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadImageHandler(t *testing.T) {
	err := initApp()
	assert.Nil(t, err)
	router := SetupRouter()

	w, err := newfileUploadRequest(router, "/api/v1/images", map[string]string{"width":"200","height":"150"}, "file", "testdata/test_image.jpg")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetResizedImagesHandler(t *testing.T) {
	err := initApp()
	assert.Nil(t, err)
	router := SetupRouter()

	w := performRequest(router, "GET", "/api/v1/images/1/resized", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResizeImageHandler(t *testing.T) {
	err := initApp()
	assert.Nil(t, err)
	router := SetupRouter()

	w := performRequest(router, "POST", "/api/v1/images/1/resized", bytes.NewBufferString(url.Values{"width":{"200"},"height":{"180"}}.Encode()))
	assert.Equal(t, http.StatusOK, w.Code)
}

//

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if method == "POST" {
	 req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func newfileUploadRequest(r http.Handler, uri string, params map[string]string, paramName, path string) (*httptest.ResponseRecorder, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}