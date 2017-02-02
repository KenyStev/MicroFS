package main

import(
	"./microfs"
	// "math"
	"fmt"
)

var mountedDisk *microfs.Disk

func main() {
	mountedDisk = nil
	for ;; {
		var command, diskname string
		var size, sizeOfBlock int64
		fmt.Scanf("%s %s %d %d",&command,&diskname,&sizeOfBlock,&size)

		if mountedDisk == nil {
			if command == "mkdisk" {
				microfs.CreateDisk(diskname,size,sizeOfBlock)
			}else if command == "mount" {
				mountedDisk = microfs.Mount(diskname,sizeOfBlock)
			}else if command == "unmount" || 
					command == "allocate" || 
					command == "unallocate"||
					command == "info"{
				fmt.Println("command '",command,"' need a mounted Disk")
			}else{
				fmt.Println("command '",command,"' not exist")
			}
		}else{
			if command == "unmount" {
				microfs.Unmount(mountedDisk)
			}else if diskname == "block" {
				if command == "allocate"{
					_ = mountedDisk.AllocateBlock()
				}else if command == "unallocate" {
					mountedDisk.UnallocateBlock(sizeOfBlock)
				}
			}else if command == "info" {
				mountedDisk.PrintInfo()
			}else if command == "mkdisk" || command == "mount" {
				fmt.Println("command '",command,"' can not be use with mounted Disk")
			}else{
				fmt.Println("command '",command,"' not exist")
			}
		}
	}
}