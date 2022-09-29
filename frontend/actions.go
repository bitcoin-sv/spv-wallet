package frontend

import (
	"embed"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	apirouter "github.com/mrz1836/go-api-router"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed react/build/*
var embeddedFiles embed.FS

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router) {

	s := apirouter.NewStack()

	// This is needed so we do not mess up with other API endpoints :-(
	webPaths := []string{
		"/images/*.",
		"/static/*.",
		"/assets-manifest.json",
		"/favicon.ico",
		"/index.html",
		"/manifest.json",
		"/robots.txt",
	}

	sub, err := fs.Sub(embeddedFiles, "react/build")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := statigz.FileServer(
		sub.(fs.ReadDirFS),
		brotli.AddEncoding,
		statigz.EncodeOnInit,
	)

	// register all the embedded files
	for _, path := range webPaths {
		router.HTTPRouter.GET(
			path,
			s.Wrap(router.Request(apirouter.StandardHandlerToHandle(func() http.HandlerFunc {
				return fileServer.ServeHTTP
			}()))),
		)
	}

	// register the root path "/"
	router.HTTPRouter.GET(
		"/",
		s.Wrap(router.Request(apirouter.StandardHandlerToHandle(staticFileHandler()))),
	)

	// register the 404 handler
	router.HTTPRouter.NotFound = staticFileHandler()
}

func staticFileHandler() http.HandlerFunc {
	return func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			staticFile, err := getFile(req)
			if err != nil {
				log.Print(err)
			}

			_, err = w.Write(staticFile)
			if err != nil {
				log.Print(err)
			}
		}
	}()
}

func getFile(req *http.Request) (indexFile []byte, err error) {
	path, _ := os.Getwd()
	if path != "" {
		filePath := filepath.Clean(req.URL.Path)
		if filePath == "/" || filePath == "" {
			filePath = "/index"
		}
		fileName := filepath.Clean(path + "/frontend/react/build/static-html" + filePath + ".html")
		if _, err = os.Stat(fileName); !errors.Is(err, os.ErrNotExist) {
			indexFile, _ = ioutil.ReadFile(fileName)
		}
	}

	// fallback to the basic index.html file - this will load the site as a normal react site (non cached)
	if indexFile == nil {
		indexFile, err = embeddedFiles.ReadFile("react/build/index.html")
		if err != nil {
			log.Print(err)
		}
	}

	// fallback if we find nothing at all
	if indexFile == nil {
		indexFile = []byte("ERROR: Cannot find index file")
	}

	return
}
