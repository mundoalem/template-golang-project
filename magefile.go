//go:build mage

// This file is part of template-golang-project.
//
// template-golang-project is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// template-golang-project is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with template-golang-project. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

const (
	CoverageMode       string = "atomic"
	SupportedPlatforms string = "darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"
)

var (
	BuildDir            string = path.Join(".", "build")
	CmdDir              string = path.Join(".", "cmd")
	CoverageProfileFile string = path.Join(".", "coverage.txt")
	DistDir             string = path.Join(".", "dist")
	InternalDir         string = path.Join(".", "internal")
	PkgDir              string = path.Join(".", "pkg")
	ReleaseDir          string = path.Join(BuildDir, "release")
	VendorDir           string = path.Join(".", "vendor")
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// commit returns the hash of the last commit done to the git repository HEAD
func commit() (string, error) {
	hash, err := sh.Output("git", "rev-parse", "--short", "HEAD")

	if err != nil {
		return "", err
	}

	return hash, nil
}

// constains returns whether a string is inside a slice or not
func contains(s []string, el string) bool {
	for _, v := range s {
		if v == el {
			return true
		}
	}

	return false
}

// linkerFlags returns a list of argument flags to be used bby the go compiler while building
// the project
func linkerFlags(isRelease bool) (string, error) {
	var err error
	var flags []string

	flags = append(flags, fmt.Sprintf(`-X "main.BuildTime=%s"`, time.Now().Format(time.RFC3339)))
	hash, err := commit()

	if err != nil {
		fmt.Printf("Warning: Could not get commit hash, reason: %s", err)
	} else {
		flags = append(flags, fmt.Sprintf(`-X "main.Commit=%s"`, hash))
	}

	versionTag, err := version()

	if err != nil {
		fmt.Printf("Warning: Could not get version tag, reason: %s\n", err)
	} else {
		flags = append(flags, fmt.Sprintf(`-X "main.Version=%s"`, versionTag))
	}

	if isRelease {
		flags = append(flags, "-s", "-w")
	}

	return strings.Join(flags, " "), nil
}

// version returns the version of the project or the last revision hash if a tag cannot be found
func version() (string, error) {
	var err error
	var revision string

	if os.Getenv("CI") != "" {
		revision = os.Getenv("GITHUB_SHA")
	} else {
		revision, err = sh.Output("git", "rev-list", "--tags", "--max-count=1")

		if err != nil {
			return "", err
		}
	}

	if revision == "" {
		return "", errors.New("Could not find repository's revision hash")
	}

	tag, err := sh.Output("git", "describe", "--always", "--tags", revision)

	if err != nil {
		return "", err
	}

	tag = strings.TrimSpace(tag)

	if tag == "" {
		return "", errors.New("Release tag not found")
	}

	return tag, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MAGE TARGETS
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Build compiles the project for all the supported platforms
func Build() error {
	pterm.Info.Println("Build process started")

	isRelease := false

	if os.Getenv("CI") != "" {
		isRelease = true
	}

	pterm.Info.Printfln("Release mode: %t", isRelease)

	var err error
	var commandDirs []string

	err = filepath.Walk(CmdDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != CmdDir {
			commandDirs = append(commandDirs, path)
		}

		return nil
	})

	if err != nil {
		pterm.Error.Printfln("Could not find command directories: ", err)
		return err
	}

	ldFlags, err := linkerFlags(isRelease)

	if err != nil {
		pterm.Error.Printfln("Could not get linker flags: ", err)
		return err
	}

	hasErrors := false

	var wg sync.WaitGroup

	for _, platform := range strings.Split(SupportedPlatforms, " ") {
		for _, commandDir := range commandDirs {
			wg.Add(1)

			go func(platform string, dir string) {
				defer wg.Done()

				command := path.Base(dir)
				os := strings.Split(platform, "/")[0]
				arch := strings.Split(platform, "/")[1]
				mainFilePath := path.Join(dir, "main.go")
				outputFilePath := path.Join(BuildDir, os+"_"+arch, command)

				envVars := map[string]string{
					"CGO_ENABLED":          "0",
					"GO111MODULE":          "on",
					"GO15VENDOREXPERIMENT": "1",
					"GOOS":                 os,
					"GOARCH":               arch,
				}

				args := []string{
					"build",
					"-o",
					outputFilePath,
					"-mod=vendor",
					"-ldflags=" + ldFlags,
					mainFilePath,
				}

				err := sh.RunWithV(envVars, "go", args...)

				if err != nil {
					pterm.Error.Printfln("Build for %s/%s failed: %s", os, arch, err)
					hasErrors = true
				} else {
					pterm.Info.Printfln("Build for %s/%s completed", os, arch)
				}
			}(platform, commandDir)
		}
	}

	wg.Wait()

	if hasErrors {
		pterm.Error.Printfln("Build process failed!")
		return errors.New("One or more build processes failed!")
	}

	pterm.Success.Println("Build process completed")

	return nil
}

// Clean removes temporary and build files
func Clean() error {
	pterm.Info.Println("Clean process started")

	pathsToRemove := []string{
		CoverageProfileFile,
	}

	err := filepath.Walk(BuildDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != BuildDir {
			pathsToRemove = append(pathsToRemove, path)
		}

		return nil
	})

	if err != nil {
		pterm.Error.Printfln("Failed to determine paths to be removed: ", err)
		return err
	}

	cleanReport := make([]pterm.BulletListItem, 0)
	flagError := false

	for _, path := range pathsToRemove {
		cleanReporItem := pterm.NewBulletListItemFromString(path, "")

		if err := sh.Rm(path); err != nil {
			flagError = true
			cleanReporItem.TextStyle = pterm.NewStyle(pterm.FgRed)
			cleanReporItem.BulletStyle = pterm.NewStyle(pterm.FgRed)

			pterm.Error.Printfln("Could not remove '%s': %s", path, err)
		}

		cleanReport = append(cleanReport, cleanReporItem)
	}

	pterm.DefaultSection.Println("Removed")
	pterm.DefaultBulletList.WithItems(cleanReport).Render()

	if flagError {
		pterm.Error.Println("Clean process completed with errors")
		return errors.New("Process failed")
	}

	pterm.Success.Println("Clean process completed")

	return nil
}

// Lint checks the project's code for style and syntax issues
func Lint() error {
	pathsToLint := []string{
		CmdDir,
		InternalDir,
		PkgDir,
	}

	for _, path := range pathsToLint {
		var args []string
		var err error

		if os.Getenv("CI") != "" {
			args = []string{
				"-e",
				"-l",
				path,
			}

			output, err := sh.Output("gofmt", args...)

			if err != nil {
				return err
			}

			if strings.TrimSpace(output) != "" {
				filesWithErrors := strings.Join(strings.Split(output, "\n"), ", ")
				errorMessage := fmt.Sprintf("Some files need linting: %s\n", filesWithErrors)

				return errors.New(errorMessage)
			}
		} else {
			args = []string{
				"fmt",
				fmt.Sprintf("./%s/...", path),
			}

			err = sh.RunV("go", args...)

			if err != nil {
				return err
			}
		}
	}

	for _, path := range pathsToLint {
		var args []string
		var err error

		args = []string{
			"-tests",
			"-f",
			"stylish",
			fmt.Sprintf("./%s/...", path),
		}

		err = sh.RunV("staticcheck", args...)

		if err != nil {
			return err
		}
	}

	return nil
}

// Lock syncs the go.sum file and the project's dependencies
func Lock() error {
	if err := sh.RunV("go", "mod", "vendor"); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

// Release generates a release tarball containing the built files for each supported platform
func Release() error {
	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	programName := path.Base(wd)
	version, err := version()

	if err != nil {
		return err
	}

	platforms := strings.Split(SupportedPlatforms, " ")

	for _, platform := range platforms {
		goos := strings.Split(platform, "/")[0]
		goarch := strings.Split(platform, "/")[1]
		sourcePath := path.Join(BuildDir, fmt.Sprintf("%s_%s", goos, goarch))

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("Warning: Build files not found for %s/%s\n", goos, goarch)
			continue
		}

		tarball := fmt.Sprintf("%s-%s-%s-%s.tar.gz", programName, version[1:], goos, goarch)
		destPath := path.Join(DistDir, tarball)

		args := []string{
			"-czf",
			destPath,
			"-C",
			sourcePath,
			".",
		}

		err = sh.RunV("tar", args...)

		if err != nil {
			return err
		}

		fmt.Printf("Release package for %s/%s created\n", goos, goarch)
	}

	if os.Getenv("CI") != "" {
		fmt.Printf("::set-output name=version::%s\n", version)
	}

	return nil
}

// Reset removes all files that Clean does plus the vendor directory
func Reset() error {
	mg.Deps(Clean)

	if err := sh.Rm(VendorDir); err != nil {
		return err
	}

	if err := os.Mkdir(VendorDir, 0755); err != nil {
		return err
	}

	return nil
}

// Scan runs a security check using Snyk to search for known vulnerabilities in project
func Scan() error {
	_, err := exec.LookPath("snyk")

	if err != nil {
		return err
	}

	args := []string{
		"test",
		"--fail-on=upgradable",
	}

	return sh.RunV("snyk", args...)
}

// Test runs the unit test for the project
func Test() error {
	args := []string{
		"test",
		"-v",
		"-count=1",
		"-race",
		"-coverprofile=" + CoverageProfileFile,
		"-covermode=" + CoverageMode,
		"./...",
	}

	return sh.RunV("go", args...)
}
