// +build mage

package main

import (
	"context"
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
)

const (
	BuildDir            string = "build"
	CmdDependencies     string = "snyk"
	CmdDir              string = "cmd"
	CoverageProfileFile string = "coverage.txt"
	CoverageMode        string = "atomic"
	InternalDir         string = "internal"
	PkgDir              string = "pkg"
	SupportedPlatforms  string = "darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"
	VendorDir           string = "vendor"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func commit() (string, error) {
	hash, err := sh.Output("git", "rev-parse", "--short", "HEAD")

	if err != nil {
		return "", err
	}

	return hash, nil
}

func contains(s []string, el string) bool {
	for _, v := range s {
		if v == el {
			return true
		}
	}

	return false
}

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

	tag, err := sh.Output("git", "describe", "--tags", revision)

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

func Build() error {
	isRelease := false

	if os.Getenv("CI") != "" {
		isRelease = true
	}

	fmt.Printf("Release mode: %t\n", isRelease)

	var err error
	var commandDirs []string

	err = filepath.Walk(CmdDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != CmdDir {
			commandDirs = append(commandDirs, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	ldFlags, err := linkerFlags(isRelease)

	if err != nil {
		fmt.Println(err)
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
				outputFilePath := path.Join(PkgDir, os+"_"+arch, command)

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
					fmt.Printf("Error: Build for %s/%s failed: %s\n", os, arch, err)
					hasErrors = true
				} else {
					fmt.Printf("Build for %s/%s completed\n", os, arch)
				}
			}(platform, commandDir)
		}
	}

	wg.Wait()

	if hasErrors {
		return errors.New("One or more build processes failed!")
	}

	return nil
}

func Clean() error {
	pathsToRemove := []string{
		CoverageProfileFile,
	}

	err := filepath.Walk(PkgDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != PkgDir {
			pathsToRemove = append(pathsToRemove, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, path := range pathsToRemove {
		if err := sh.Rm(path); err != nil {
			return err
		}
	}

	return nil
}

func Lint() error {
	pathsToLint := []string{
		CmdDir,
		InternalDir,
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

	return nil
}

func Lock() error {
	if err := sh.RunV("go", "mod", "vendor"); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

func Release(ctx context.Context) error {
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
		sourcePath := path.Join(PkgDir, fmt.Sprintf("%s_%s", goos, goarch))

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("Warning: Build files not found for %s/%s\n", goos, goarch)
			continue
		}

		tarball := fmt.Sprintf("%s-%s-%s-%s.tar.gz", programName, version[1:], goos, goarch)
		destPath := path.Join(BuildDir, tarball)

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
