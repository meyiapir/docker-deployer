package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

//func checkFilePath(filePath string) bool {
//	files, err := ioutil.ReadDir(filePath)
//	if err != nil {
//		log.Fatal(err)
//		return false
//	}
//	for _, file := range files {
//		if file.Name() == "Dockerfile" {
//			return true
//		}
//	}
//	return false
//}

func checkProject(projectName string) string {
	var a string
	projects := ReadFile(ConfigDir + "projects.json")
	for range projects {
		if getKey(projects, projectName, "name") != "false" {
			a = "true"
		} else {
			a = "false"
		}
	}
	return a
}

func buildImage(imageName string, pathToProject string) bool {
	fmt.Println("Путь к проекту: " + pathToProject)
	cmd := exec.Command("docker", "build", pathToProject, "-t", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return true
}

func tagImage(imageName string) bool {
	cmd := exec.Command("docker", "tag", imageName, registry+"/"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return true
}

func pushImage(imageName string) bool {
	cmd := exec.Command("docker", "push", registry+"/"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return true
}

func deleteImage(imageName string) bool {
	var images [2]string
	images[0] = imageName
	images[1] = registry + "/" + imageName
	for image := range images {
		fmt.Println("Удаление образа " + images[image])
		cmd := exec.Command("docker", "rmi", images[image])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return true
}

func btpImage(imageName string, projectName string, pathToProject string) bool {
	buildImage(imageName, pathToProject)
	fmt.Println(" ")
	fmt.Println("Проект " + projectName + " собран\n" + "Имя образа: " + imageName)

	tagImage(imageName)
	fmt.Println(" ")
	fmt.Println("Образ " + imageName + " отмечен как " + registry + "/" + imageName + "\n")

	pushImage(imageName)
	fmt.Println(" ")
	fmt.Println("Образ " + imageName + " отправлен в репозиторий " + registry + "\n")
	return true
}

func main() {
	var projectName string
	var imageName string
	var isDelete bool
	fmt.Print("Введите имя проекта: ")
	fmt.Fscan(os.Stdin, &projectName)
	if checkProject(projectName) == "true" {
		fmt.Println("Проект " + projectName + " найден\n")
		fmt.Println("**********************************************************")
		fmt.Println("Введите имя образа: ")
		fmt.Fscan(os.Stdin, &imageName, "\n")
		fmt.Println(" ")
		start := time.Now()
		btpImage(imageName, projectName, getKey(ReadFile(ConfigDir+"projects.json"), projectName, "path"))
		if getKey(ReadFile(ConfigDir+"config.json"), "delete_after_push", "enabled") == "true" {
			deleteImage(imageName)
			fmt.Println(" ")
			fmt.Println("Образ " + imageName + " удалён\n")
			isDelete = true
		} else {
			fmt.Println(" ")
			fmt.Println("Образ " + imageName + " не удалён, так как выключена соответствующая настройка \n")
		}
		end := time.Now()
		fmt.Println(" ")
		fmt.Println("Время выполнения: " + fmt.Sprint(end.Sub(start)))
		if isDelete {
			fmt.Println("Название образа: " + imageName + " (был удалён)")
		} else {
			fmt.Println("Название образа: " + imageName)
		}
		fmt.Println("Проект: " + projectName)

	} else {
		fmt.Println("Проект " + projectName + " не найден")
	}

	time.Sleep(43200 * time.Second)
}
