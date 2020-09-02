package ops

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/feiyuw/simgo/logger"
	"github.com/labstack/echo/v4"
)

const (
	UPLOADED_DIR = "./upload"
)

func uploadFile(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := ioutil.TempFile(UPLOADED_DIR, "*."+file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"filepath": dst.Name()})
}

func removeFile(c echo.Context) error {
	// avoid delete unexpected file
	filePath := filepath.Join(UPLOADED_DIR, filepath.Base(c.QueryParam("filepath")))
	logger.Warn("ops/files", filePath)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
