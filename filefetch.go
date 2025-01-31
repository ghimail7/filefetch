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

// Init Flags
// General Flags
var longMode = flag.BoolP("long", "l", false, "Use Long Mode.")
var dirFirst = flag.BoolP("directoriesfirst", "d", false, "List Directories before Files.")
var humanReadable = flag.BoolP("humanreadable", "h", true, "Simplifies sized to abbreviated binary units.")

// Long Mode Flags
var sectionMargin = flag.IntP("sectionmargin", "s", 1, "How far the sections/file names will be")
var lastModifiedEnable = flag.BoolP("lastmodified", "m", false, "Enable the Last Modified Section on Long Mode.")
var permsEnable = flag.BoolP("permissions", "p", false, "Enable the Perms Section on Long Mode.")

// Format Flags
var dateFormat = flag.StringP("format", "f", "02/01/2006 15:04:05.000", "Date `format` for Last Modified, if enabled.")

// Init Functions
func sep(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

func bft(num int64) string {
	returns := sep(num)
	if *humanReadable {
		var shortened float64 = float64(num);
		prefixes := []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi", "Yi"}
		for _, pref := range prefixes {
			if shortened < 1024 {
				returns = decimal.NewFromFloat(shortened).RoundBank(2).String() + pref
				break
			}
			shortened /= 1024
			
		}
	}

	return returns
}

// Print w/Ansi & Truncation if true
func PrintPA(bol bool, str string, trl int, ansi string) {
	fmt.Print(ansi + str + "\033[0m")
	if trl > 0 {
		fmt.Print(strings.Repeat(" ", trl - len(str)))
	}
}

func printFiles(dirOnly bool) {
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		if dirOnly && file.IsDir() || !dirOnly && !file.IsDir() {
			// If the file is a directory, use blue.
			shortAnsi := "\033[34m"
			fileSuf := ""
			if !dirOnly {
				shortAnsi = "\033[32m" 
				fileSuf = "\n"
			}

			// Long Mode: line-separated file information.
			// Short Mode: space-separated file names.
			if *longMode {
				// Directories will only display directory in the size section.
				if dirOnly { 
					PrintPA(true, "Directory", hl.size, "\033[36m") 
				} else { 
					PrintPA(true, bft(fileinfo.Size()) + "B", hl.size, "\033[36m") 
				} 

				// Whether to enable the Last Modified or Permissions sections.
				PrintPA(*lastModifiedEnable, fileinfo.ModTime().Format(*dateFormat), hl.modi, "\033[33m")
				if dirOnly { 
					PrintPA(*lastModifiedEnable, strings.Replace(fileinfo.Mode().Perm().String(), "-", "d", 1), hl.perm, "\033[31m") 
				} else {
					PrintPA(*lastModifiedEnable, fileinfo.Mode().Perm().String(), hl.perm, "\033[31m") 
				} 
			}

			PrintPA(true, file.Name(), 0, shortAnsi)
			fmt.Print(fileSuf)
		}
	}	
}

// Main Loop
func main() {
	flag.Parse()
	// Checks lengths of each sections for proper alignment.
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fileinfo, _ := os.Stat(file.Name())
		sizeString := bft(fileinfo.Size()) + "B"
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
		PrintPA(true, "Size", hl.size, "\033[1;4m")
		PrintPA(*lastModifiedEnable, "Last Modified", hl.modi, "\033[1;4m") 
		PrintPA(*permsEnable, "Permissions", hl.perm, "\033[1;4m")
		PrintPA(true, "Name", hl.size, "\033[1;4m")
		fmt.Println("\033[0m")
	}

	if *dirFirst {
		printFiles(true)	
		printFiles(false)	
	} else {
		printFiles(false)	
		printFiles(true)	
	}

	plural := "ies"
	if dircount == 1 { plural = "y" }
	if !*longMode { fmt.Println() }
	fmt.Println("Fetched \033[36;1m" + bft(int64(dirsize)) + "B \033[0mof Files and \033[34;1m" + sep(int64(dircount)) + " \033[0mDirector" + plural + ".")
}
