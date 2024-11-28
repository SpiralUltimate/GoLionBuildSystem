package buildSystem

import "os"

// Instantiate a builder class
type Builder struct {
	// Project name
	project string
	// Compiler
	compiler string
	// Project files
	projectFiles []string
}

// Parses Builder struct into a cmake file, using the given name for the resulting cmake file
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

// Makes a new project, using the given name
func (build *Builder) Project(name string, compiler string) {
	build.project = name
	build.compiler = compiler
}

// Adds files to the project
func (build *Builder) AddFiles(files ...string) {
	build.projectFiles = append(build.projectFiles, files...)
}
