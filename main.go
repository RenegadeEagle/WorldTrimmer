/* Internally minecraft stores region files in .mca format
* This project is used to delete region files outside a search range. Perfect for when your worlds have become too big and you need to trim them down.
* The convention for storing region files is r.x.z.mca (Region, X Unit, Z Unit)
* To calculate the blocks the region handles, multiple the X and Z by 512.
* Those two numbers are the start coordinates for the region.
* Add 511 to both of those values to find the end coordinates. You then have a 512x512 area for the region.
* Example:
* r.4.-5.mca
* 4 * 512 = 2048, -5 * 512 = -2,560
* Add 511 to find the end point for the region to get the coords:
* 2559, -2049
* You now know the exact blocks that are stored in this .mca file.
* We use this information to delete .mca files that are outside a specified region.
* To keep worlds smaller for Minecraft servers if you forgot to set a world border at one point.
*/

package main

import "fmt"
import "os"
import "strconv"
import "time"
import "io/ioutil"
import "path/filepath"
import "strings"
import "math"

//   ./wc limit debug path
//   ./wc 25000 true /root/myserver/myworld/region
func main() {
    if len(os.Args) < 4 {
        fmt.Println("Not enough arguments!")
        os.Exit(0)
    } else {
        fmt.Println("Starting WorldTrimmer v1.0!")
        debug, _ := strconv.ParseBool(os.Args[2])
        deleteThreshhold, _ := strconv.Atoi(os.Args[1])
        path := os.Args[3]

        if debug {
            fmt.Println("Starting WorldTrimmer in DEBUG mode! In debug mode, the files will not actually be deleted.")
            fmt.Println("Delete Threshold is " + os.Args[1] + ", path for world files is \"" + path + "\".")
        }

        fmt.Println("WARNING! We are going to begin scanning the folder " + path + " and deleting any files past the threshold you have provided.")
        fmt.Println("There is NO way to recover your deleted world past this point. ")
        fmt.Println("You have ten seconds to Control-C this process before deletion begins.")
        time.Sleep(10 * time.Second)
        fmt.Println("Beginning to scan files. This could take a while.")

        files, _ := ioutil.ReadDir(path)
        fmt.Printf("Scanning %d files!\n", len(files))
        counter := 0
        deletedFiles := 0
        var bytes int64
        for _, f := range files {
            counter++

            fmt.Printf("[%d/%d] Scanning File %s\n", counter, len(files), f.Name())

            var ext = filepath.Ext(path + f.Name())

            if ext != ".mca"  {
                continue //Don't try this on non mca files. That would be bad.
            }

            split := strings.Split(f.Name(), ".")

            x, _ := strconv.Atoi(split[1])
            z, _ := strconv.Atoi(split[2])

            startX := x * 512
            endX := startX + 511

            startZ := z * 512
            endZ := startZ + 511

            absoluteX := math.Abs(float64(startX)) //New to golang. This is some strange casting...
            absoluteZ := math.Abs(float64(startZ))

            absoluteEndX := math.Abs(float64(endX))
            absoluteEndZ := math.Abs(float64(endZ))

            if x >= 0 && z >= 0 {
                if int(absoluteX) > deleteThreshhold || int(absoluteZ) > deleteThreshhold {
                    fmt.Printf("Deleting file %s at coordinates at (%d, %d)!\n", f.Name(), startX, startZ)
                    deletedFiles++
                    bytes += f.Size()
                    deleteFile(path, f.Name(), debug)
                }
                continue
            }
            //Lets calculate the closer point for negative numbers.
            if x < 0 && z < 0 {
                if int(absoluteEndX) > deleteThreshhold || int(absoluteEndZ) > deleteThreshhold {
                    fmt.Printf("Deleting file %s at coordinates at (%d, %d)!\n", f.Name(), startX, startZ)
                    deletedFiles++
                    bytes += f.Size()
                    deleteFile(path, f.Name(), debug)
                }
                continue
            }

            if int(absoluteX) > deleteThreshhold && int(absoluteEndX) > deleteThreshhold {
                fmt.Printf("Deleting file %s at coordinates at (%d, %d)!\n", f.Name(), startX, startZ)
                deletedFiles++
                bytes += f.Size()
                deleteFile(path, f.Name(), debug)
                continue
            }
            if int(absoluteZ) > deleteThreshhold && int(absoluteEndZ) > deleteThreshhold {
                fmt.Printf("Deleting file %s at coordinates at (%d, %d)!\n", f.Name(), startX, startZ)
                deletedFiles++
                bytes += f.Size()
                deleteFile(path, f.Name(), debug)
                continue
            }
        }
        var mb float64
        mb = float64(bytes) / float64(1000000)
        fmt.Printf("\n======== Scan Complete ========\n")
        fmt.Printf("%d Files Scanned, %d files deleted.\n", counter, deletedFiles)
        fmt.Printf("Deleted %f megabytes worth of files!\n", mb)
        fmt.Printf("===============================\n")
    }
}


func deleteFile(path string, name string, debug bool) {
    if debug {
        return
    }
    err := os.Remove(path + "/" + name)
    if err != nil {
        fmt.Println("Error deleting file!", err)
    }
}
