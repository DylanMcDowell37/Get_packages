package main

import (
	"fmt"
	"strings"
	"github.com/anaskhan96/soup"
	"net/http"
	"os"
	"io"
	"bufio"
	"log"
	"flag"
)

func souper(package_url string) []string{
	resp, err := soup.Get(package_url)
	if err != nil {
		os.Exit(1)
	}
	var releases []string
	doc := soup.HTMLParse(resp)
	links := doc.Find("table").FindAll("a")
	for _, link := range links {
		releases = append(releases, link.Attrs()["href"])	
	}
	return releases
}

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func contains(arr []string, str string) bool {
    for _, item := range arr {
        if item == str {
            return true
        }
    }
    return false
}

func fileread(file_name string) []string {
	var output_array []string 
	file, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open or create file: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		package_name := scanner.Text()
		output_array = append(output_array, package_name)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return output_array
}

func appendToFile(fileName, text string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(text + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func createDirectory(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func getpackage(maindir string){
	package_name := fileread("Packages.txt")
	for i := 0; i < len(package_name); i++ {
		replacer := strings.NewReplacer(" patches", "", " updates", "")
		name := replacer.Replace(package_name[i])
		history := fileread(fmt.Sprintf("links/%s.txt", name))
		package_name_url := strings.Replace(package_name[i], " ", "-", -1)
		package_url := fmt.Sprintf("https://www.manageengine.com/products/desktop-central/patch-management/%s.html", package_name_url)
		if contains(history, package_url){
			fmt.Printf("%s is up to date\n", name)
		}else {
			output := souper(package_url)
			new_url := output[len(output)-2]
			downloads := souper(fmt.Sprintf("https://www.manageengine.com/products/desktop-central/patch-management/%s", new_url))
			lastSlashIndex := strings.LastIndex(downloads[len(downloads)-1], "/")
			name_out := downloads[len(downloads)-1][lastSlashIndex+1:]
			appendToFile(fmt.Sprintf("links/%s.txt", name), downloads[len(downloads)-1])
			createDirectory(fmt.Sprintf("%s/%s", maindir, name))
			out := fmt.Sprintf("%s/%s/%s", maindir, name, name_out)
			fmt.Printf("%s update is downloading\n", name)
			DownloadFile(downloads[len(downloads)-1], out)
		}
	}
}

func main(){
	maindir := flag.String("dir", "", "Directory path to create")
	flag.Parse()
	if *maindir == "" {
		fmt.Println("Please provide a directory path using -dir")
		return
	}
	
	// Ensure main directory exists
	if err := createDirectory(*maindir); err != nil {
		log.Fatalf("Failed to create main directory: %v", err)
	}
	
	getpackage(*maindir)
}