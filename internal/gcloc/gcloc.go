package gcloc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FSFile struct {
	Path      string
	Content   string
	Count     uint
	Extension string
}

type GClocImpl struct {
	files      *[]FSFile
	totalCount uint
}

func NewGClocImpl() *GClocImpl {
	return &GClocImpl{}
}

func (g *GClocImpl) Run() {
	files, err := getAllFiles()
	if err != nil {
		panic(err)
	}
	processedFiles, totalCount, err := processFiles(*files)
	if err != nil {
		panic(err)
	}
	g.files = processedFiles
	g.totalCount = totalCount

	fmt.Println(g.totalCount)
}

func getAllFiles() (*[]FSFile, error) {

	files := []FSFile{}

	autoExclusionList := []string{".git", "license", ".md", "node_modules"}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			found := false
			for _, autoExclusionValue := range autoExclusionList {
				if strings.Contains(strings.ToLower(path), strings.ToLower(autoExclusionValue)) {
					found = true
				}
			}
			if !found {
				files = append(files, FSFile{Path: strings.ToLower(path), Extension: filepath.Ext(path)})
			}
		}

		return nil
	})

	return &files, err
}

func processFiles(files []FSFile) (*[]FSFile, uint, error) {

	totalCount := 0

	for index, fsFile := range files {
		if fsFile.Extension != "" && (fsFile.Extension == ".go" || fsFile.Extension == ".java") {
			file, err := os.Open(fsFile.Path)
			if err != nil {
				log.Fatalf("failed to open file: %s", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineCount := 0

			for scanner.Scan() {
				lineCount++
			}

			if err := scanner.Err(); err != nil {
				log.Fatalf("error reading file: %s", err)
			}

			files[index].Count = uint(lineCount)
			totalCount += lineCount
		}
	}

	return &files, uint(totalCount), nil
}
