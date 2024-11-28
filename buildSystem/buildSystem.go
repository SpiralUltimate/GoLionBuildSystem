package buildSystem

import (
	"os"
	"os/exec"

	LionFormat "github.com/SpiralUltimate/GoLionFormat/format"
)

// Instantiate a builder class
type Builder struct {
	// Project name
	project string
	// Compiler
	compiler string
	// C++ standard
	cppStandard int
	// config to run built MakeFiles in (eg. "debug", "release", etc.)
	config string

	// Project files
	projectFiles []string
}

// Makes a new project, using the given name. Also sets the compiler (eg. "g++", "gcc")
func (build *Builder) Project(name string, compiler string) {
	build.project = name
	build.compiler = compiler
}

// Sets a C++ standard
func (build *Builder) CppStandard(cppStandard int) {
	build.cppStandard = cppStandard
}

// Sets the config to use (eg. "debug", "release")
func (build *Builder) Config(config string) {
	build.config = config
}

// Adds files to the project
func (build *Builder) Files(files ...string) {
	build.projectFiles = append(build.projectFiles, files...)
}

// Parses Builder struct into a cmake file, using the given path for the resulting cmake file
func (build *Builder) Parse(cmakeFilePath string) error {
	// Create a new CMake file
	cmakeFile, err := os.Create(cmakeFilePath)
	// Error check cmake file
	if err != nil {
		return err
	}
	// Close cmake file at end of function
	defer cmakeFile.Close()

	// Write project name and compiler info
	cmakeFile.WriteString("cmake_minimum_required(VERSION 3.16)\n")
	cmakeFile.WriteString("project(" + build.project + ")\n")
	cmakeFile.WriteString("set(CMAKE_CXX_STANDARD 17)\n")
	cmakeFile.WriteString("set(CMAKE_CXX_COMPILER " + build.compiler + ")\n\n")

	// Write project files
	for _, file := range build.projectFiles {
		cmakeFile.WriteString("add_executable(" + build.project + " " + file + ")\n")
	}

	// Return nil as error to indicate success
	return nil
}

// Runs an already parsed cmake file
func (build *Builder) Run() error {
	// Build cmake into MakeFiles inside a "build" directory
	buildCmd := exec.Command("cmake", "-S", ".", "-B", "build", "-G", "\"Unix Makefiles\"")
	// Set standard input and output to current terminal stdin and stdout
	buildCmd.Stdin = os.Stdin
	buildCmd.Stdout = os.Stdout
	// Run build command
	err := buildCmd.Run()
	// Check for build command errors
	if err != nil {
		return err
	}

	// Get the config param string
	config, err := LionFormat.Format("config={s}", build.config)
	// Check for errors when formatting
	if err != nil {
		return err
	}
	// Run built MakeFiles
	buildMakeCmd := exec.Command("make", config)
	// Set standard input and output to current terminal stdin and stdout
	buildMakeCmd.Stdin = os.Stdin
	buildMakeCmd.Stdout = os.Stdout
	// Run buildMake command
	err = buildMakeCmd.Run()
	// Check for buildMake command errors
	if err != nil {
		return err
	}

	// Return nil as error to indicate success
	return nil
}
