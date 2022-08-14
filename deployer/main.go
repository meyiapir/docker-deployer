package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var BaseDir = getBaseDir()
var ConfigDir = BaseDir + "\\configs\\"
var registry = getKey(ReadFile(ConfigDir+"config.json"), "registry", "address") +
	":" + getKey(ReadFile(ConfigDir+"config.json"), "registry", "port")

func getBaseDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// Чтение json и сохранение вложенных данных в ключи

func ReadFile(fileName string) map[string]interface{} {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		print("File error: %v\n", err)
	}
	var data map[string]interface{}
	json.Unmarshal(file, &data)
	return data
}

func getKey(data map[string]interface{}, key1 string, key2 string) string {
	gettingKey, err := data[key1].(map[string]interface{})
	if err != true {
		return "false"
	}
	copyKey := gettingKey[key2]
	str := fmt.Sprintf("%v", copyKey)
	return str
}

func checkProject(projectName string) bool {
	var a bool
	projects := ReadFile(ConfigDir + "projects.json")
	for range projects {
		if getKey(projects, projectName, "name") != "false" {
			a = true
		} else {
			a = false
		}
	}
	return a
}

func buildImage(imageName string, projectName string) bool {
	pathToProject := getKey(ReadFile(ConfigDir+"projects.json"), projectName, "path")
	cmd := exec.Command("docker", "build", pathToProject, "-t", imageName)
	//cmd := exec.Command("calc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println(" ")
	fmt.Println("Проект " + projectName + " собран\n" + "Имя образа: " + imageName)
	return true
}

func tagImage(imageName string) bool {
	cmd := exec.Command("docker", "tag", imageName, registry+"/"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println(" ")
	fmt.Println("Образ " + imageName + " отмечен как " + registry + "/" + imageName + "\n")
	return true

}

func pushImage(imageName string) bool {
	cmd := exec.Command("docker", "push", registry+"/"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println(" ")
	fmt.Println("Образ " + imageName + " отправлен в репозиторий " + registry + "\n")
	return true
}

func main() {
	var projectName string
	var imageName string
	fmt.Print("Введите имя проекта или путь к Dockerfile: ")
	fmt.Fscan(os.Stdin, &projectName)
	if checkProject(projectName) {
		fmt.Println("Проект " + projectName + " найден\n")
		fmt.Println("**********************************************************")
		fmt.Println("Введите имя образа: ")
		fmt.Fscan(os.Stdin, &imageName, "\n")
		fmt.Println(" ")
		buildImage(imageName, projectName)
		tagImage(imageName)
		//pushImage(imageName)
	} else {
		fmt.Println("Проект " + projectName + " не найден")
	}
	//time.Sleep(3600 * time.Second)
}

//cmd := exec.Command("dir", "C:\\Users\\meyap\\go\\src\\Sportifier-deployer")
//getKey(ReadFile(ConfigDir+"projects.json"), "mailer", "path")

//func main1() {
//	if getKey(ReadFile(ConfigDir+"projects.json"), "Sender", "name") != false {
//		fmt.Println("yes")
//	} else {
//		fmt.Println("no")
//	}
//}
