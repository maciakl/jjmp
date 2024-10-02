package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"

	"github.com/charmbracelet/huh"
    "github.com/charmbracelet/lipgloss"
)

const version = "0.1.0"

var db map[string]string

func main() {

    // set ansi profile to retain color when piping output
    lipgloss.SetColorProfile(0)

    bookmarks := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

    // load the db file
    LoadDB()
    defer SaveDB()

    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "-v", "--version":
            Version()
            os.Exit(0)
        case "-h", "--help":
            Usage()
            os.Exit(0) 
        case "set":
            if len(os.Args) != 3 {
                Usage()
                os.Exit(1)
            } else {
                // check if second arg is in the bookmarks list
                if slices.Contains(bookmarks, os.Args[2]) {
                    // get current working directory
                    cwd, err := os.Getwd()
                    if err != nil {
                        fmt.Fprintln(os.Stderr, "Error:", err)
                        os.Exit(1)
                    }
                    db[os.Args[2]] = cwd
                    SaveDB()
                } else {
                    fmt.Fprintln(os.Stderr, "Invalid bookmark")
                    Usage()
                    os.Exit(1)
                }
            }

        case "delete":
            if len(os.Args) != 3 {
                fmt.Fprintln(os.Stderr, "Invalid number of arguments")
                Usage()
                os.Exit(1)
            } else {
                if slices.Contains(bookmarks, os.Args[2]) {
                    deleteBookmark(os.Args[2])
                    //delete(db, os.Args[2])
                    //SaveDB()
                } else {
                    fmt.Fprintln(os.Stderr, "Invalid bookmark")
                    Usage()
                    os.Exit(1)
                }
            }

        case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
            fmt.Println(db[os.Args[1]])

        default:
            Usage()
            os.Exit(1)
        } 
    } else {
        // if map is empty, print error message
        if len(db) == 0 {
            fmt.Fprintln(os.Stderr, "No bookmarks found")
            os.Exit(1)
        } else {
            path := ShowSelect()
            fmt.Println(path)
        }
    }
    
}

func Version() {
    fmt.Fprintln(os.Stderr, filepath.Base(os.Args[0]), "version", version)
}

func Usage() {
    fmt.Fprintln(os.Stderr, "Usage:", filepath.Base(os.Args[0]), "[options]")
    fmt.Fprintln(os.Stderr, "Options:")
    fmt.Fprintln(os.Stderr, "  -v, --version    Print version information and exit")
    fmt.Fprintln(os.Stderr, "  -h, --help       Print this message and exit")
    fmt.Fprintln(os.Stderr, "  set <0-9>        Bookmark the current directory")
    fmt.Fprintln(os.Stderr, "  delete <0-9>     Delete a bookmark")
    fmt.Fprintln(os.Stderr, "  <0-9>            Change to the bookmarked directory")
}

// function that reads ~/.jjmpdb text file and loads it into the db map
func LoadDB() {
    db = make(map[string]string)

    homedir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

    dbpath := path.Join(homedir, ".jjmpdb")

    file, err := os.Open(dbpath)
    if err != nil {
        SaveDB()
    }
    defer file.Close()

    // deserialize the db map from the text filepath
    dec := gob.NewDecoder(file)
    err = dec.Decode(&db)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

}

// function that writes the db map to the ~/.jjmpdb text file overwriting the existing content
func SaveDB() {

    homedir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

    dbpath := path.Join(homedir, ".jjmpdb")

    file, err := os.Create(dbpath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }
    defer file.Close()

    // serialize the db map to the text file overwriting the existing content
    enc := gob.NewEncoder(file)
    err = enc.Encode(db)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }
}

func ShowSelect() string {

    var mypath string

    choices := make([]string, 0, len(db))
    for key, value := range db {
        choices = append(choices, key + "\t" + value)
    }

    sort.Strings(choices)

    selector:=huh.NewSelect[string]().Title("Choose a bookmark:").
    Options(huh.NewOptions(choices...)...).
    Value(&mypath)
    err := selector.Run()
    if err != nil {
        os.Exit(1)
    }

    // split on tab to get the path
    mypath = mypath[len(mypath)-len(mypath)+2:]
    return mypath

}

func deleteBookmark(key string) {
    value := db[key]

    var result bool
    fmt.Fprintln(os.Stderr, "Deleting bookmark:", key, "\t", value, "?")
    c:= huh.NewConfirm()
    c.Value(&result).Title("Are you sure?").Run()

    if result {
        delete(db, key)
        SaveDB()
    }
}
