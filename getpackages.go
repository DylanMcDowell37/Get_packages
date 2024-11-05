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
	//Get url output from soup
	resp, err := soup.Get(package_url)
	if err != nil {
		os.Exit(1)
	}
	//Create empty array
	var releases []string
	//Parse soup data
	doc := soup.HTMLParse(resp)
	//Find all links located in a table
	links := doc.Find("table").FindAll("a")
	//Append all links into array
	for _, link := range links {
		releases = append(releases, link.Attrs()["href"])	
	}
	return releases
}

func DownloadFile(url string, filepath string) error {
	//Create file path
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	//http get request for data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//Inject(download) data into newly created file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func contains(arr []string, str string) bool {
	//Check if string exist in array
    for _, item := range arr {
        if item == str {
            return true
        }
    }
    return false
}

func fileread(file_name string) []string {
	//Set array variable
	var output_array []string 
	//Open file
	file, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open or create file: %v\n", err)
	}
	defer file.Close()
	//create scan of file
	scanner := bufio.NewScanner(file)
	//loop through scan and append data into array
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
	//Open file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()
	//Append text to file and create newline
	if _, err := file.WriteString(text + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func createDirectory(dirPath string) error {
	//Create directory if it does not exist
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func getpackage(maindir string){
	//Get list of packages to download
	package_name := fileread("Packages.txt")
	for i := 0; i < len(package_name); i++ {
		//Set variables
		replacer := strings.NewReplacer(" patches", "", " updates", "")
		name := replacer.Replace(package_name[i])
		history := fileread(fmt.Sprintf("links/%s.txt", name))
		package_name_url := strings.Replace(package_name[i], " ", "-", -1)
		package_url := fmt.Sprintf("https://www.manageengine.com/products/desktop-central/patch-management/%s.html", package_name_url)
		//Get output from first url
		output := souper(package_url)
		//New url
		new_url := output[len(output)-2]
		//Check history
		if contains(history, new_url){
			fmt.Printf("%s is up to date\n", name)
		}else {
			//using new url get link for downloading latest package
			downloads := souper(fmt.Sprintf("https://www.manageengine.com/products/desktop-central/patch-management/%s", new_url))
			//Set variables
			lastSlashIndex := strings.LastIndex(downloads[len(downloads)-1], "/")
			name_out := downloads[len(downloads)-1][lastSlashIndex+1:]
			//Add to history file
			appendToFile(fmt.Sprintf("links/%s.txt", name), new_url)
			//Create download directory if it does not exist in maindir
			createDirectory(fmt.Sprintf("%s/%s", maindir, name))
			//Set variables and Download package
			out := fmt.Sprintf("%s/%s/%s", maindir, name, name_out)
			fmt.Printf("%s update is downloading\n", name)
			DownloadFile(downloads[len(downloads)-1], out)
		}
	}
}

func main(){
	//Get main dir from -dir flag
	maindir := flag.String("dir", "", "Directory path to create")
	flag.Parse()
	//Throw error if -dir was never specified
	if *maindir == "" {
		fmt.Println("Please provide a directory path using -dir")
		return
	}
	
	// Ensure main directory exists
	if err := createDirectory(*maindir); err != nil {
		log.Fatalf("Failed to create main directory: %v", err)
	}

	// Ensure links directory exists
	if err := createDirectory("links"); err != nil {
		log.Fatalf("Failed to create links directory: %v", err)
	}
	//Run getpackage
	getpackage(*maindir)
}