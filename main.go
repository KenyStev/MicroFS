package main

import(
	"./microfs"
	// "math"
	"fmt"
)

// var mountedDisk *Disk

func main() {
	for ;; {
		var command, diskname string
		var size, sizeOfBlock int64
		fmt.Scanf("%s %s %d %d",&command,&diskname,&size,&sizeOfBlock)

		if command == "mkdisk" {
			microfs.CreateDisk(diskname,size,sizeOfBlock)
		}
	}
}