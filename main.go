package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	srcDir   = "src"
	libDir   = "libs"
	buildDir = "build"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("[ERROR] Missing arguments.\n[HELP] Usage: jargo <command>")
		return
	}

	command := args[1]
	switch command {
	case "new":
		projectName := args[2]
		newProject(projectName)
	case "run":
		projectName := "."
		if len(args) > 2 {
			projectName = args[2]
		}
		runProject(projectName)
	default:
		fmt.Println("[ERROR] Invalid command.\n[INFO] Commands: new, run")
	}
}

func newProject(projectDirectory string) {
	projectPath := filepath.Join(projectDirectory)

	err := os.Mkdir(projectDirectory, 0755)
	if err != nil {
		fmt.Println("[ERROR] Directory already exists.")
		return
	}

	err = os.Mkdir(filepath.Join(projectPath, srcDir), 0755)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	err = os.Mkdir(filepath.Join(projectPath, libDir), 0755)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// create starter boilerplate in src
	file, err := os.Create(fmt.Sprintf("%s/src/Main.java", projectDirectory))
	file.WriteString(MainBoilerplate)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	file2, err := os.Create(fmt.Sprintf("%s/src/Test.java", projectDirectory))
	file2.WriteString(TestBoilerplate)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file2.Close()

	err = os.WriteFile(projectDirectory+"/libs/Example.jar", JarData, 0644)
	if err != nil {
		panic(err)
	}
}

func runProject(projectDirectory string) {
	projectPath := filepath.Join(projectDirectory)
	buildPath := filepath.Join(projectPath, buildDir)

	// create build dir if it doesnt exist else clear build folder
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		err = os.Mkdir(buildPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
	} else {
		dirEntries, err := os.ReadDir(buildPath)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		// Iterate over the directory entries
		for _, entry := range dirEntries {
			// Check if the entry is a file and ends with ".class"
			if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".class") {
				// Build the file path
				filePath := buildDir + entry.Name()
				os.Remove(filePath)
			}
		}
	}

	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}

	compileCommand := fmt.Sprintf("javac -cp \"%[1]s/libs/*\" -d %[1]s/build %[1]s/src/*.java", projectDirectory)
	executeCommand := fmt.Sprintf("java -cp \"%[1]s/build%s%[1]s/libs/*\" Main", projectDirectory, separator)
	runCommand(compileCommand)
	runCommand(executeCommand)
}

func runCommand(command string) {
	debug := os.Getenv("DEBUG")
	if debug == "true" || debug == "1" {
		fmt.Println(command)
	}
	cmd := exec.Command("bash", "-c", command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	}
	// link the executed command to the current terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
