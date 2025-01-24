package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/cornfeedhobo/pflag"
)

// sep formats a number with commas for thousands separators.
func sep(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

// PrintPA prints a string with ANSI color codes and padding.
func PrintPA(str string, pad int, ansi string) {
	if pad > 0 {
		fmt.Print(ansi + str + "\033[0m" + strings.Repeat(" ", pad-len(str)))
	} else {
		fmt.Print(ansi + str + "\033[0m")
	}
}

// printFiles prints the details of files in the current directory based on the provided flags.
func printFiles(dirOnly bool, longMode bool, lastModifiedEnable bool, permsEnable bool, dateFormat string, hlsize, hlmodi, hlperms int) {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fileinfo, err := os.Stat(filepath.Join(".", file.Name()))
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		if (dirOnly && file.IsDir()) || (!dirOnly && !file.IsDir()) {
			if longMode {
				if dirOnly {
					PrintPA("Directory", hlsize, "\033[36m")
				} else {
					PrintPA(sep(fileinfo.Size())+"B", hlsize, "\033[36m")
				}
				if lastModifiedEnable {
					PrintPA(fileinfo.ModTime().Format(dateFormat), hlmodi, "\033[33m")
				}
				if permsEnable {
					if dirOnly {
						PrintPA(strings.Replace(fileinfo.Mode().Perm().String(), "-", "d", 1), hlperms, "\033[31m")
					} else {
						PrintPA(fileinfo.Mode().Perm().String(), hlperms, "\033[31m")
					}
				}
				if dirOnly {
					PrintPA(file.Name(), 0, "\033[34m")
				} else {
					PrintPA(file.Name(), 0, "\033[32m")
				}
				fmt.Println("")
			} else {
				if dirOnly {
					PrintPA(file.Name(), 0, "\033[34m")
					fmt.Print(" ")
				} else if !file.IsDir() {
					PrintPA(file.Name(), 0, "\033[32m")
					fmt.Print(" ")
				}
			}
		}
	}
}

func main() {
	longMode := flag.BoolP("long", "l", false, "Use Long Mode.")
	dirFirst := flag.BoolP("directoriesfirst", "d", false, "List Directories before Files.")
	lastModifiedEnable := flag.BoolP("lastmodified", "m", false, "Enable the Last Modified Section on Long Mode.")
	permsEnable := flag.BoolP("permissions", "p", false, "Enable the Perms Section on Long Mode.")
	dateFormat := flag.StringP("format", "f", "02/01/2006 15:04:05.000", "Date `format` for Last Modified, if enabled.")
	flag.Parse()

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var hlsize, hlmodi, hlperms int = 4, 13, 11
	var dirsize int64
	var dircount int

	// Determine the maximum lengths for size, last modified, and permissions columns for proper alignment.
	for _, file := range files {
		fileinfo, err := os.Stat(filepath.Join(".", file.Name()))
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		if len(sep(fileinfo.Size())+"B") > hlsize {
			hlsize = len(sep(fileinfo.Size()) + "B")
		}
		if len(fileinfo.ModTime().Format(*dateFormat)) > hlmodi {
			hlmodi = len(fileinfo.ModTime().Format(*dateFormat))
		}
		if len(fileinfo.Mode().Perm().String()) > hlperms {
			hlperms = len(fileinfo.Mode().Perm().String())
		}

		if file.IsDir() {
			dircount++
		} else {
			dirsize += fileinfo.Size()
		}
	}

	hlsize++
	hlmodi++
	hlperms++

	if *longMode {
		// Print headers for the long mode.
		PrintPA("Size", hlsize, "\033[1;4m")
		if *lastModifiedEnable {
			PrintPA("Last Modified", hlmodi, "\033[1;4m")
		}
		if *permsEnable {
			PrintPA("Permissions", hlperms, "\033[1;4m")
		}
		PrintPA("Name", hlsize, "\033[1;4m")
		fmt.Println("\033[0m")
	}

	// Print directories and files based on the specified flags.
	if *dirFirst {
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat, hlsize, hlmodi, hlperms)
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat, hlsize, hlmodi, hlperms)
	} else {
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat, hlsize, hlmodi, hlperms)
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *dateFormat, hlsize, hlmodi, hlperms)
	}

	plural := "ies"
	if dircount == 1 {
		plural = "y"
	}
	if !*longMode {
		fmt.Println()
	}
	// Print summary of fetched directories and file sizes.
	fmt.Println("Fetched \033[36;1m" + sep(dirsize) + "B \033[0mof Files and \033[34;1m" + sep(int64(dircount)) + " \033[0mDirector" + plural + ".")
}
