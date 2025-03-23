# WinPkg - Advanced Windows Package Manager

WinPkg is an **unofficial**, **advanced** package manager for Windows, designed to simplify the installation, management, and updating of software packages on Windows systems. It is written in Go and provides a command-line interface (CLI) for users to interact with a repository of packages, enabling seamless installation, uninstallation, searching, and updating of software.

This tool is ideal for developers, system administrators, and power users who want to automate software management on Windows systems. WinPkg is lightweight, fast, and extensible, allowing users to create and publish their own packages to the repository.

## Features

- **Package Installation**: Install software packages with a single command.
- **Package Uninstallation**: Remove installed packages cleanly and efficiently.
- **Package Search**: Search for packages in the repository by name.
- **Package Listing**: View all installed packages and their versions.
- **Package Update**: Update installed packages to their latest versions.
- **Custom Commands**: Execute custom commands directly from the CLI.
- **Package Publishing**: Create and publish your own packages to the repository.
- **Repository Management**: Load and manage a repository of packages locally.
- **Multi-Package Manager Support**: Install and uninstall packages using popular package managers like `pip`, `nuget`, `npm`, `go`, `gem`, and `composer`.
- **Importing from GitHub**: Install and Uninstall Packages from GitHub 

## Installation

To install WinPkg, follow these steps:

1. **Download the Executable**:
   - Download the latest version of `winpkg6.exe` from the official GitHub repository:
     [Download winpkg6.exe](https://github.com/CalestialAshley35/WinPkg/blob/main/winpkg6.exe)

2. **Run the Executable**:
   - Place the downloaded `winpkg5.exe` file in a directory of your choice.
   - Open a Command Prompt or PowerShell window and navigate to the directory containing `winpkg6.exe`.
   - Run the executable to start using WinPkg:
     ```bash
     ./winpkg6.exe
     ```

3. **Verify Installation**:
   - Once the application starts, you should see the `winpkg>` prompt, indicating that WinPkg is ready to use.

Or you can download from https://calestialashley35.itch.io/winpkg

## Usage

WinPkg provides a simple and intuitive CLI for managing packages. Below are the available commands and their usage:

### 1. Install a Package
To install a package, use the `install` command followed by the package name. You can optionally specify a version or a package manager flag and also importing from github.

```bash
winpkg> install <package_name> [version] [flag]
```

Example:
```bash
winpkg> install notepad++
winpkg> install requests -python
```

Example with github:
```bash
winpkg> install <package-name> -github
```

### 2. Uninstall a Package
To uninstall a package, use the `uninstall` command followed by the package name. You can optionally specify a package manager flag.

```bash
winpkg> uninstall <package_name> [flag]
```

Example:
```bash
winpkg> uninstall notepad++
winpkg> uninstall requests -python
```

### 3. List Installed Packages
To list all installed packages, use the `list` command.

```bash
winpkg> list
```

### 4. Search for a Package
To search for a package in the repository, use the `search` command followed by the package name.

```bash
winpkg> search <package_name>
```

Example:
```bash
winpkg> search python
```

### 5. Update a Package
To update an installed package to the latest version, use the `update` command followed by the package name.

```bash
winpkg> update <package_name>
```

Example:
```bash
winpkg> update python
```

### 6. Execute a Custom Command
To execute a custom command, use the `use` command followed by the command.

```bash
winpkg> use <command>
```

Example:
```bash
winpkg> use dir
```

### 7. Publish a Package
To publish a new package to the repository, use the `publish` command followed by the path to the `winpkg.infoi` file.

```bash
winpkg> publish <path_to_winpkg.infoi>
```

Example:
```bash
winpkg> publish C:\packages\my_package.winpkg.infoi
```

### 8. Exit the Application
To exit the WinPkg CLI, use the `exit` command.

```bash
winpkg> exit
```

## Creating and Publishing Packages

WinPkg allows users to create and publish their own packages to the repository. To create a package, follow these steps:

1. **Define Package Details**:
   - Use the `createPackage` function in the CLI to define the package details, including:
     - Name
     - Version
     - Installation file
     - Description
     - Install command
     - Section

2. **Save Package Information**:
   - The package details are saved in a `winpkg.infoi` file, which can be published to the repository.

3. **Publish the Package**:
   - Use the `publish` command to add the package to the repository.

Example `winpkg.infoi` file:
```
*Name*: my_package
*Version*: 1.0.0
*Installation*: my_package.exe
*Description*: A sample package for WinPkg
*Install*: my_package.exe /S
*Section*: Utilities
```

## Repository Management

WinPkg maintains a local repository of packages. The repository is loaded when the application starts. If no packages are found, users are prompted to create a new package.

### Loading the Repository
The repository is loaded automatically when WinPkg starts. If no packages are found, users can create a new package using the `createPackage` function.

### Adding Packages to the Repository
Packages can be added to the repository by publishing a `winpkg.infoi` file using the `publish` command.

## Technical Details

### Package Structure
Each package in the repository is represented by a `Package` struct with the following fields:
- **Name**: The name of the package.
- **Version**: The version of the package.
- **Installation**: The installation file or command.
- **Description**: A brief description of the package.
- **InstallCmd**: The command used to install the package.
- **Section**: The category or section of the package.

### Installation Process
When a package is installed, WinPkg executes the `InstallCmd` specified in the package details. The package is then added to the `installedPackages` map.

### Uninstallation Process
When a package is uninstalled, WinPkg uses the appropriate package manager command (e.g., `pip`, `npm`, `msiexec`) to remove the package from the system.

## License

WinPkg is released under the MIT License. See the [LICENSE](https://github.com/CalestialAshley35/WinPkg/blob/main/LICENSE) file for more details.

## Disclaimer

WinPkg is an **unofficial** package manager and is not affiliated with any official Windows software distribution channels. However, it is a trusted tool developed with transparency and security in mind. Use it at your own risk. The developers are not responsible for any issues arising from the use of this tool.

## Support

For support or to report issues, please open an issue on the [GitHub repository](https://github.com/CalestialAshley35/WinPkg).
