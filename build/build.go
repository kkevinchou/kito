package main

import (
	"fmt"
	"os/exec"
)

func main() {
	mainFolder := "build"
	buildFolder := "build_output"
	shaderFolder := "shaders"
	assetsFolder := "_assets"

	cmd := exec.Command("rm", "-r", buildFolder)
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}

	cmd = exec.Command("mkdir", buildFolder)
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}

	cmd = exec.Command("cp", "-r", shaderFolder, buildFolder)
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}

	cmd = exec.Command("cp", "-r", assetsFolder, buildFolder)
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}

	cmd = exec.Command("cp", "-r", mainFolder+"/SDL2.dll", buildFolder)
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}

	cmd = exec.Command("go", "build", "-o", buildFolder+"/kito.exe")
	if output, err := cmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(string(output))
	}
}
