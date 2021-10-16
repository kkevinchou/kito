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
	configFile := "config.json"

	cmd := exec.Command("rm", "-r", buildFolder)
	runCMD(cmd)

	cmd = exec.Command("mkdir", buildFolder)
	runCMD(cmd)

	cmd = exec.Command("cp", configFile, buildFolder)
	runCMD(cmd)

	cmd = exec.Command("cp", "-r", shaderFolder, buildFolder)
	runCMD(cmd)

	cmd = exec.Command("cp", "-r", assetsFolder, buildFolder)
	runCMD(cmd)

	cmd = exec.Command("cp", "-r", mainFolder+"/SDL2.dll", buildFolder)
	runCMD(cmd)

	cmd = exec.Command("go", "build", "-o", buildFolder+"/kito.exe")
	runCMD(cmd)
}

func runCMD(cmd *exec.Cmd) {
	if _, err := cmd.Output(); err != nil {
		fmt.Printf("failed to run command %s\n%s\n", cmd.Path, err)
	} else {
		fmt.Printf("- %v\n", cmd.Args)
	}
}
