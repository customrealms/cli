package papermc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

type cachedFetcher struct {
	Fetcher  Fetcher
	cacheDir string
}

func NewCachedFetcher(fetcher Fetcher) (Fetcher, error) {

	// Setup the cache directory
	cacheDir, _ := os.UserCacheDir()
	cacheDir = path.Join(cacheDir, "cr-cli-cache")
	if err := os.MkdirAll(cacheDir, 0777); err != nil {
		return nil, err
	}

	// Create the cached fetcher instance
	return &cachedFetcher{
		Fetcher:  fetcher,
		cacheDir: cacheDir,
	}, nil

}

func (f *cachedFetcher) getJarCacheFilename(version *Version) string {
	return path.Join(f.cacheDir, fmt.Sprintf("paper-%s-%d.jar", version.Version, version.Build))
}

func (f *cachedFetcher) findJarFile(version *Version) (io.ReadCloser, error) {

	// Get the filename of the JAR cache
	jarCacheFilename := f.getJarCacheFilename(version)

	// Check if the file exists
	stat, err := os.Stat(jarCacheFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if stat.IsDir() {
		return nil, errors.New("jar cache location is a directory")
	}

	// Read the file
	return os.Open(jarCacheFilename)

}

func (f *cachedFetcher) storeJarFile(reader io.Reader, version *Version) (string, error) {

	// Get the filename of the JAR cache
	jarCacheFilename := f.getJarCacheFilename(version)

	// Create the file for the cache
	cacheFile, err := os.Create(jarCacheFilename)
	if err != nil {
		return "", err
	}
	defer cacheFile.Close()

	// Copy the jar data to the file
	if _, err := io.Copy(cacheFile, reader); err != nil {
		return "", err
	}

	// Return the jar filename
	return jarCacheFilename, nil

}

func (f *cachedFetcher) Fetch(version *Version) (io.ReadCloser, error) {

	// Check for the file in the cache, and return the cached version is there is one
	jarReader, err := f.findJarFile(version)
	if err != nil {
		return nil, err
	}
	if jarReader != nil {
		return jarReader, nil
	}

	// Fetch the JAR file from the upstream fetcher
	res, err := f.Fetcher.Fetch(version)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// Store the jar file contents
	jarFilename, err := f.storeJarFile(res, version)
	if err != nil {
		return nil, err
	}

	// Open and return the cache file
	return os.Open(jarFilename)

}
