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

var installedPackages = make(map[string]string)
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
			flag := ""
			if len(args) > 2 {
				if strings.HasPrefix(args[2], "-") {
					flag = args[2]
				} else {
					version = args[2]
				}
			}
			installPackage(args[1], version, flag)
		case "uninstall":
			if len(args) < 2 {
				fmt.Println("Please provide a package name.")
				continue
			}
			flag := ""
			if len(args) > 2 {
				flag = args[2]
			}
			uninstallPackage(args[1], flag)
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
	if len(repository) == 0 {
		fmt.Println("No packages found. Would you like to create a new package? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(response) == "y" {
			createPackage()
		}
	}
}

func createPackage() {
	fmt.Println("Enter the package details:")

	var pkg Package
	fmt.Print("Name: ")
	fmt.Scanln(&pkg.Name)
	fmt.Print("Version: ")
	fmt.Scanln(&pkg.Version)
	fmt.Print("Installation File: ")
	fmt.Scanln(&pkg.Installation)
	fmt.Print("Description: ")
	fmt.Scanln(&pkg.Description)
	fmt.Print("Install Command: ")
	fmt.Scanln(&pkg.InstallCmd)
	fmt.Print("Section: ")
	fmt.Scanln(&pkg.Section)

	repository = append(repository, pkg)
	fmt.Println("Package created successfully!")

	savePackageInfoToFile(pkg)
}

func savePackageInfoToFile(pkg Package) {
	fileContent := fmt.Sprintf(`*Name*: %s
*Version*: %s
*Installation*: %s
*Description*: %s
*Install*: %s
*Section*: %s
`, pkg.Name, pkg.Version, pkg.Installation, pkg.Description, pkg.InstallCmd, pkg.Section)

	err := ioutil.WriteFile("winpkg.infoi", []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error saving winpkg.infoi file:", err)
	} else {
		fmt.Println("Package information saved to winpkg.infoi.")
	}
}

func installPackage(packageName, version, flag string) {
	packageInfo, err := findPackage(packageName, version)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if _, ok := installedPackages[packageName]; ok {
		fmt.Println("Package already installed:", packageName)
		return
	}

	var cmd *exec.Cmd
	switch flag {
	case "-python":
		cmd = exec.Command("pip", "install", packageInfo.Name)
	case "-nuget":
		cmd = exec.Command("nuget", "install", packageInfo.Name)
	case "-github":
		// Install from GitHub repository (example: https://github.com/user/repo)
		fmt.Println("Installing package from GitHub:", packageInfo.Name)
		cmd = exec.Command("git", "clone", packageInfo.Installation) // Assuming `Installation` is a GitHub URL
		if err := cmd.Run(); err != nil {
			fmt.Println("Error during GitHub installation:", err)
			return
		}
		// Optionally, run any setup commands post-clone, if necessary
		fmt.Println("Package installed from GitHub.")
	default:
		fmt.Println("Installing", packageInfo.Name, "version", packageInfo.Version)
		cmd = exec.Command(packageInfo.InstallCmd)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error during installation:", err)
		return
	}

	installedPackages[packageName] = packageInfo.Version
	fmt.Println("Package", packageInfo.Name, "installed successfully.")
}

func uninstallPackage(packageName, flag string) {
	if _, ok := installedPackages[packageName]; !ok {
		fmt.Println("Package not installed:", packageName)
		return
	}

	var cmd *exec.Cmd
	switch flag {
	case "-python":
		cmd = exec.Command("pip", "uninstall", "-y", packageName)
	case "-nuget":
		cmd = exec.Command("nuget", "uninstall", packageName)
	default:
		cmd = exec.Command("msiexec", "/x", packageName, "/quiet")
	}

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

	uninstallPackage(packageName, "")
	installPackage(packageName, latestPackage.Version, "")
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