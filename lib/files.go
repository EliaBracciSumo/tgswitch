package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// RenameFile : rename file name
func RenameFile(src string, dest string) {
	err := os.Rename(src, dest)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// RemoveFiles : remove file
func RemoveFiles(src string) {
	files, err := filepath.Glob(src)
	if err != nil {

		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

// CheckFileExist : check if file exist in directory
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

//CreateDirIfNotExist : create directory if directory does not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Creating directory for terragrunt: %v", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Unable to create directory for terragrunt: %v", dir)
			panic(err)
		}
	}
}

//WriteLines : writes into file
func WriteLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		return err
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	return nil
}

// ReadLines : Read a whole file into the memory and store it as array of lines
func ReadLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

//IsDirEmpty : check if directory is empty (TODO UNIT TEST)
func IsDirEmpty(name string) bool {

	exist := false

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		exist = true
	}
	return exist // Either not empty or error, suits both cases
}

//CheckDirHasTGBin : // check binary exist (TODO UNIT TEST)
func CheckDirHasTGBin(dir, prefix string) bool {

	exist := false

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
		//return exist, err
	}
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
			res = append(res, filepath.Join(dir, f.Name()))
			exist = true
		}
	}
	return exist
}

//CheckDirExist : check if directory exist
//dir=path to file
//return path to directory
func CheckDirExist(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// Path : returns path of directory
// value=path to file
func Path(value string) string {
	return filepath.Dir(value)
}

func GetCurrentDirectory() string {

	dir, err := os.Getwd() //get current directory
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}
	return dir
}

// FileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetHomeDirectory : return the home directory
func GetHomeDirectory() string {

	homedir, errHome := homedir.Dir()
	if errHome != nil {
		log.Printf("Failed to get home directory %v\n", errHome)
		os.Exit(1)
	}

	return homedir
}

// GetFileName : remove file ext.  .tgswitch.config returns .tgswitch
func GetFileName(configfile string) string {
	return strings.TrimSuffix(configfile, filepath.Ext(configfile))
}

// Print message reading file content of :
func ReadingFileMsg(filename string) {
	fmt.Printf("Reading file %s \n", filename)
}

//retrive file content of regular file
func RetrieveFileContents(file string) string {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error: %s\nFailed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/tgswitch/blob/master/README.md\n", err, file)
	}
	tgversion := strings.TrimSuffix(string(fileContents), "\n")
	return tgversion
}

// Print invalid TG version
func PrintInvalidTGVersion() {
	fmt.Println("Version does not exist or invalid terragrunt version format.\n Format should be #.#.# or #.#.#-@# where # are numbers and @ are word characters.\n For example, 0.11.7 and 0.11.9-beta1 are valid versions")
}
