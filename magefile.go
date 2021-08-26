// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		return strings.Join(flags, " "), err
	}

	flags = append(flags, fmt.Sprintf(`-X "main.Commit=%s"`, hash))

	versionTag, err := version()

	if err != nil {
		return "", err
	}

	flags = append(flags, fmt.Sprintf(`-X "main.Version=%s"`, versionTag))

	if isRelease {
		flags = append(flags, "-s", "-w")
	}

	return strings.Join(flags, " "), nil
}

func version() (string, error) {
	output, err := sh.Output("git", "tag", "--sort=-version:refname", "-l", "v*")

	if err != nil {
		return "", err
	}

	if strings.TrimSpace(output) == "" {
		return "", errors.New("Release tag not found")
	}

	tags := strings.Split(output, "\n")

	if len(tags) <= 0 {
		return "unknown", nil
	}

	return tags[0], nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// MAGE TARGETS
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func Build(ctx context.Context) error {
	mg.SerialDeps(Check, Clean, Lock)

	isRelease := false

	if v := ctx.Value("isRelease"); v != nil {
		if value, ok := v.(bool); ok {
			isRelease = value
		}
	}

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
					fmt.Printf("Build for %s/%s failed: %s\n", os, arch, err)
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

func Check() error {
	for _, cmd := range strings.Split(CmdDependencies, " ") {
		_, err := exec.LookPath(cmd)

		if err != nil {
			log.Println(err)
			return err
		}
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
		PkgDir,
	}

	for _, path := range pathsToLint {
		path = fmt.Sprintf("./%s/...", path)

		if err := sh.RunV("go", "fmt", path); err != nil {
			return err
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
	mg.CtxDeps(context.WithValue(ctx, "isRelease", true), Build)

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
		os := strings.Split(platform, "/")[0]
		arch := strings.Split(platform, "/")[1]
		sourcePath := path.Join(PkgDir, fmt.Sprintf("%s_%s", os, arch))
		tarball := fmt.Sprintf("%s-%s-%s-%s.tar.gz", programName, version[1:], os, arch)
		destPath := path.Join(BuildDir, tarball)

		args := []string{
			"-czvf",
			destPath,
			"-C",
			sourcePath,
			".",
		}

		err = sh.RunV("tar", args...)

		if err != nil {
			return err
		}
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
