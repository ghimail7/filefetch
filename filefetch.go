package main

import (
	"fmt"
	"os"
	"strings"
	"golang.org/x/text/message"
	"golang.org/x/text/language"
	"github.com/shopspring/decimal"
)

import flag "github.com/cornfeedhobo/pflag"

// Init Structs
type sec struct {
	size int
	modi int
	perm int
}

// Init Variables
var hl = sec{
	size: 4,
	modi: 13,
	perm: 11,
}
var dircount = 0
var dirsize int64 = 0

// Init Functions
func sep(num int64, humanReadable bool) string {
	returns := ""
	if humanReadable {
		var shortened float64 = float64(num);
		prefixes := []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi", "Yi"}
		for _, pref := range prefixes {
			if shortened < 1024 {
				returns = decimal.NewFromFloat(shortened).RoundBank(2).String() + pref
				break
			} else {
				shortened /= 1024
			}
		}
	} else {
		p := message.NewPrinter(language.English)
		returns = p.Sprintf("%d", num)
	}

	return returns
}

func PrintPA(str string, pad int, ansi string) {
	if pad > 0 {
		fmt.Print(ansi + str + "\033[0m" + strings.Repeat(" ", pad - len(str)))
	} else {
		fmt.Print(ansi + str + "\033[0m")
	}
}

func printFiles(dirOnly bool, longMode bool, lastModifiedEnable bool, permsEnable bool, humanReadable bool, dateFormat string) {
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if dirOnly && file.IsDir() || !dirOnly && !file.IsDir() {
			// If the file is a directory, use blue.
			shortAnsi := "\033[34m"
			if !dirOnly { shortAnsi = "\033[32m" }

			// Long Mode: line-separated file information.
			// Short Mode: space-separated file names.
			if longMode {
				// Directories will only display directory in the size section.
				if dirOnly { 
					PrintPA("Directory", hl.size, "\033[36m") 
				} else { 
					PrintPA(sep(fileinfo.Size(), humanReadable) + "B", hl.size, "\033[36m") 
				} 

				// Whether to enable the Last Modified or Permissions sections.
				if lastModifiedEnable { PrintPA(fileinfo.ModTime().Format(dateFormat), hl.modi, "\033[33m") }
				if permsEnable { 
					if dirOnly { 
						PrintPA(strings.Replace(fileinfo.Mode().Perm().String(), "-", "d", 1), hl.perm, "\033[31m") 
					} else {
						PrintPA(fileinfo.Mode().Perm().String(), hl.perm, "\033[31m") 
					} 
				}	 
				PrintPA(file.Name(), 0, shortAnsi) 
				fmt.Println("")
			} else {
				PrintPA(file.Name(), 0, shortAnsi)
				fmt.Print(" ")
			}
		}
	}	
}

// Main Loop
func main() {
	// Init Flags
	longMode := flag.BoolP("long", "l", false, "Use Long Mode.")
	dirFirst := flag.BoolP("directoriesfirst", "d", false, "List Directories before Files.")
	humanReadable := flag.BoolP("humanreadable", "h", true, "Simplifies sized to abbreviated binary units.")
	sectionMargin := flag.IntP("sectionmargin", "s", 1, "How far the sections/file names will be")

	lastModifiedEnable := flag.BoolP("lastmodified", "m", false, "Enable the Last Modified Section on Long Mode.")
	permsEnable := flag.BoolP("permissions", "p", false, "Enable the Perms Section on Long Mode.")
	
	dateFormat := flag.StringP("format", "f", "02/01/2006 15:04:05.000", "Date `format` for Last Modified, if enabled.")

	flag.Parse()

	// Checks lengths of each sections for proper alignment.
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		sizeString := sep(fileinfo.Size(), *humanReadable) + "B"
		modiString := fileinfo.ModTime().Format(*dateFormat)
		permString := fileinfo.Mode().Perm().String()

		if file.IsDir() { 
			dircount += 1 
			sizeString = "Directory"
		} else {
			dirsize += fileinfo.Size() 
		} 

		if len(sizeString) > hl.size { hl.size = len(sizeString) }
		if len(modiString) > hl.modi { hl.modi = len(modiString) }
		if len(permString) > hl.perm { hl.perm = len(permString) }
	}

	// Section Margins
	hl.size+=*sectionMargin
	hl.modi+=*sectionMargin
	hl.perm+=*sectionMargin

	if *longMode {
		PrintPA("Size", hl.size, "\033[1;4m")
		if *lastModifiedEnable { PrintPA("Last Modified", hl.modi, "\033[1;4m") }
		if *permsEnable { PrintPA("Permissions", hl.perm, "\033[1;4m") }
		PrintPA("Name", hl.size, "\033[1;4m")
		fmt.Println("\033[0m")
	}

	if *dirFirst {
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *humanReadable, *dateFormat)	
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *humanReadable, *dateFormat)	
	} else {
		printFiles(false, *longMode, *lastModifiedEnable, *permsEnable, *humanReadable, *dateFormat)	
		printFiles(true, *longMode, *lastModifiedEnable, *permsEnable, *humanReadable, *dateFormat)	
	}

	plural := "ies"
	if dircount == 1 { plural = "y" }
	if !*longMode { fmt.Println() }
	fmt.Println("Fetched \033[36;1m" + sep(int64(dirsize), *humanReadable) + "B \033[0mof Files and \033[34;1m" + sep(int64(dircount), false)+ " \033[0mDirector" + plural + ".")
}
