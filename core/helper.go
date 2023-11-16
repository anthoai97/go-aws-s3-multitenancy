package core

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetEnvVar[T any](key string, defaultValue T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	var ret any
	switch any(defaultValue).(type) {
	case string:
		ret = value

	case int:
		// don't actually ignore errors
		i, _ := strconv.ParseInt(value, 10, 64)
		ret = int(i)
	}
	return ret.(T)
}

func GetObjectExtention(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

func FormatBucketPrefixForTree(path string) string {

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	path = strings.TrimPrefix(path, "/")

	return path
}

func GenerateMetaHashObject(prefix, path string) string {
	val := prefix + FormatBucketPrefixForTree(path)

	md5 := md5.Sum([]byte(val))

	return fmt.Sprintf("%x", md5)
}

func MustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
