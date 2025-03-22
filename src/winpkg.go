package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Installation string `json:"installation"`
	Description string `json:"description"`
	InstallCmd  string `json:"install"`
	Section     string `json:"section"`
}

var installedPackages = make(map[string]string) // map of package names to versions
var repository []Package

func main() {
	loadRepository()

	for {
		fmt.Print("winpkg> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		args := strings.Fields(input)
		command := args[0]

		switch command {
		case "install":
			if len(args) < 2 {
				fmt.Println("Please provide a package name.")
				continue
			}
			version := ""
			if len(args) > 2 {
				version = args[2]
			}
			installPackage(args[1], version)
		case "uninstall":
			if len(args) < 2 {
				fmt.Println("Please provide a package name.")
				continue
			}
			uninstallPackage(args[1])
		case "list":
			listPackages()
		case "search":
			if len(args) < 2 {
				fmt.Println("Please provide a package name to search.")
				continue
			}
			searchPackage(args[1])
		case "update":
			if len(args) < 2 {
				fmt.Println("Please provide a package name to update.")
				continue
			}
			updatePackage(args[1])
		case "use":
			if len(args) < 2 {
				fmt.Println("Please provide a command to use.")
				continue
			}
			useCommand(args[1])
		case "publish":
			if len(args) < 2 {
				fmt.Println("Please provide a winpkg.infoi file.")
				continue
			}
			publishPackage(args[1])
		case "exit":
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}

func loadRepository() {
	// Ideally, we would load a repository from a JSON or database
	// For simplicity, we'll hardcode some packages here
	repository = []Package{
		{Name: "examplePackage", Version: "1.0.0", Installation: "example_v1.exe", InstallCmd: "start example_v1.exe", Section: "utilities", Description: "An example package."},
		{Name: "examplePackage", Version: "2.0.0", Installation: "example_v2.exe", InstallCmd: "start example_v2.exe", Section: "utilities", Description: "An updated example package."},
	}

	// In a real-world scenario, you could fetch this repository from a remote source
}

func installPackage(packageName, version string) {
	packageInfo, err := findPackage(packageName, version)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if _, ok := installedPackages[packageName]; ok {
		fmt.Println("Package already installed:", packageName)
		return
	}

	fmt.Println("Installing", packageInfo.Name, "version", packageInfo.Version)
	cmd := exec.Command(packageInfo.InstallCmd)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error during installation:", err)
		return
	}

	installedPackages[packageName] = packageInfo.Version
	fmt.Println("Package", packageInfo.Name, "installed successfully.")
}

func uninstallPackage(packageName string) {
	if _, ok := installedPackages[packageName]; !ok {
		fmt.Println("Package not installed:", packageName)
		return
	}

	// Assuming the package has a standard uninstall process (e.g., msiexec for MSI)
	cmd := exec.Command("msiexec", "/x", packageName, "/quiet")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error during uninstallation:", err)
		return
	}

	delete(installedPackages, packageName)
	fmt.Println("Package", packageName, "uninstalled successfully.")
}

func listPackages() {
	if len(installedPackages) == 0 {
		fmt.Println("No packages installed.")
		return
	}

	fmt.Println("Installed Packages:")
	for pkg, version := range installedPackages {
		fmt.Println(pkg, version)
	}
}

func searchPackage(packageName string) {
	var foundPackages []Package
	for _, pkg := range repository {
		if strings.Contains(pkg.Name, packageName) {
			foundPackages = append(foundPackages, pkg)
		}
	}

	if len(foundPackages) == 0 {
		fmt.Println("No packages found for", packageName)
		return
	}

	fmt.Println("Search Results:")
	for _, pkg := range foundPackages {
		fmt.Printf("%s (%s) - %s\n", pkg.Name, pkg.Version, pkg.Description)
	}
}

func updatePackage(packageName string) {
	if _, ok := installedPackages[packageName]; !ok {
		fmt.Println("Package not installed:", packageName)
		return
	}

	latestPackage, err := findPackage(packageName, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	currentVersion := installedPackages[packageName]
	if latestPackage.Version == currentVersion {
		fmt.Println("Already using the latest version of", packageName)
		return
	}

	// Uninstall the current version and install the latest version
	uninstallPackage(packageName)
	installPackage(packageName, latestPackage.Version)
}

func useCommand(command string) {
	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error during command execution:", err)
	}
}

func publishPackage(infoFile string) {
	fileContent, err := ioutil.ReadFile(infoFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	packageInfo := parsePackageInfo(fileContent)

	// Here you would add the package to a database or repository
	repository = append(repository, packageInfo)
	fmt.Println("Published Package:", packageInfo.Name, "version", packageInfo.Version)
}

func findPackage(packageName, version string) (Package, error) {
	for _, pkg := range repository {
		if pkg.Name == packageName && (version == "" || pkg.Version == version) {
			return pkg, nil
		}
	}
	return Package{}, fmt.Errorf("package not found: %s", packageName)
}

func parsePackageInfo(content []byte) Package {
	lines := strings.Split(string(content), "\n")
	packageInfo := Package{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*Name*:") {
			packageInfo.Name = strings.TrimPrefix(line, "*Name*: ")
		} else if strings.HasPrefix(line, "*Installation*:") {
			packageInfo.Installation = strings.TrimPrefix(line, "*Installation*: ")
		} else if strings.HasPrefix(line, "*Description*:") {
			packageInfo.Description = strings.TrimPrefix(line, "*Description*: ")
		} else if strings.HasPrefix(line, "*Install*:") {
			packageInfo.InstallCmd = strings.TrimPrefix(line, "*Install*: ")
		} else if strings.HasPrefix(line, "*Section*:") {
			packageInfo.Section = strings.TrimPrefix(line, "*Section*: ")
		} else if strings.HasPrefix(line, "*Version*:") {
			packageInfo.Version = strings.TrimPrefix(line, "*Version*: ")
		}
	}

	return packageInfo
}
