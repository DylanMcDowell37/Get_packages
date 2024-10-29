# Application Version Tracker

## Overview

The Application Version Tracker is a Go script designed to download the latest versions of applications from ManageEngine's Patch Management service. It maintains a history of previously downloaded versions, ensuring that you can keep track of updates efficiently.

## Features

- Downloads the latest versions of specified applications.
- Maintains a history of downloaded versions for easy reference.
- Creates directories for each application to organize downloads.
- Simple configuration via command-line arguments.

## Requirements

- Go (version 1.15 or later)
- Internet connection
- Access to ManageEngine's Patch Management URLs

## Installation

1. **Clone the repository** (if applicable):
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Install dependencies**:
   This script uses the `soup` library for HTML parsing. Install it with:
   ```bash
   go get github.com/anaskhan96/soup
   ```

3. **Ensure you have a `Packages.txt` file**:
   - Create a `Packages.txt` file in the same directory as the script.
   - List the applications you want to track, one per line (e.g., `PuTTY patches`, `Notepad++ updates`).

## Usage

Run the script from the command line, specifying the directory where downloaded applications should be stored using the `-dir` flag:

```bash
go run .\test.go -dir=C:/path/to/download/directory
```

### Command-Line Arguments

- `-dir`: The directory path where downloaded applications will be stored. This path will be created if it doesn't exist.

## Example

1. Create a `Packages.txt` file with the following content:
   ```
   PuTTY patches
   Notepad++ updates
   ```

2. Run the script:
   ```bash
   go run .\test.go -dir=C:/Users/<username>/Scripts/Go/test
   ```

3. The script will download the latest versions of the applications listed in `Packages.txt` and store them in `C:/Users/<username>/Scripts/Go/test`.

## Directory Structure

- Each application will have its own directory within the specified path, containing:
  - The latest version of the application.
  - A `links` subdirectory with historical version links.


## Acknowledgments

- This project utilizes the [soup](https://github.com/anaskhan96/soup) library for HTML parsing.
- Special thanks to the ManageEngine team for providing a robust Patch Management solution.
