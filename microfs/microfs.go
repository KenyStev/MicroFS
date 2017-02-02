package microfs

import(
	"math"
	"unsafe"
	"encoding/binary"
	"../basicfs"
	"fmt"
)

type Disk struct{
	volumeManager 	*basicfs.Volume
	size 			int64
	freeSpace 		int64
	sizeOfBlock 	int64
	blocks 			int64
	freeBlocks 		int64
	headBlock 		int64
	tailBlock 		int64
}

func (d *Disk) setSize(s int64) {
	d.size = s
}

func (d *Disk) setSizeOfBlock(s int64) {
	d.sizeOfBlock = s
}

func (d *Disk) toBytesArray() []byte {
	size := int64(unsafe.Sizeof(*d)) - 8
	buffer := make([]byte,size)

	// fmt.Printf("%v",d)

	binary.PutVarint(buffer[0:8], d.size)
	binary.PutVarint(buffer[8:16], d.freeSpace)
	binary.PutVarint(buffer[16:24], d.sizeOfBlock)
	binary.PutVarint(buffer[24:32], d.blocks)
	binary.PutVarint(buffer[32:40], d.freeBlocks)
	binary.PutVarint(buffer[40:48], d.headBlock)
	binary.PutVarint(buffer[48:56], d.tailBlock)

	// fmt.Printf("%v",toStructFS(buffer))

	return buffer
}

func toStructFS(fs *Disk, bytes []byte){
	// fs := new(Disk)
	fs.size, _ = binary.Varint(bytes[0:8])
	fs.freeSpace, _ = binary.Varint(bytes[8:16])
	fs.sizeOfBlock, _ = binary.Varint(bytes[16:24])
	fs.blocks, _ = binary.Varint(bytes[24:32])
	fs.freeBlocks,_ = binary.Varint(bytes[32:40])
	fs.headBlock,_ = binary.Varint(bytes[40:48])
	fs.tailBlock,_ = binary.Varint(bytes[48:56])

	// return fs
}

func (d *Disk) writeFS() {
	// fmt.Println(d.toBytesArray())
	buffer := make([]byte,d.sizeOfBlock)
	copy(buffer,d.toBytesArray())

	d.volumeManager.WriteBlock(0,buffer)
}

func round(size int64) int64 {
	return int64(math.Pow(2,math.Ceil(math.Log2(float64(size)))))
}

func initializeFreeBlocksList(disk *Disk) {
	buffer := make([]byte,disk.sizeOfBlock)
	for i := 0; i < int(disk.blocks); i++ {
		nextBlock := make([]byte, 8)
		binary.PutVarint(nextBlock, int64(i+1))
		copy(buffer[disk.sizeOfBlock - 8:],nextBlock)
		disk.volumeManager.WriteBlock(int64(i),buffer)
	}
	disk.headBlock = 0
	disk.tailBlock = int64(disk.blocks - 1)
	nullBlock := make([]byte, 8)
	binary.PutVarint(nullBlock, int64(-1))
	disk.volumeManager.WriteBlock(disk.tailBlock,nullBlock)
}

func Format(disk *Disk) {
	disk.freeSpace = disk.size
	disk.freeBlocks = disk.blocks
	initializeFreeBlocksList(disk)
	_ = disk.AllocateBlock()
	disk.writeFS()
}

func (d *Disk) AllocateBlock() int64 {
	var allocatedBlock int64
	if d.headBlock != -1{
		allocatedBlock = d.headBlock
		buffer := make([]byte,d.sizeOfBlock)
		d.volumeManager.ReadBlock(int64(d.headBlock),buffer)
		nextHead,_ := binary.Varint(buffer[d.sizeOfBlock - 8:])
		d.headBlock = nextHead

		d.freeBlocks -= 1
		d.freeSpace -= d.sizeOfBlock

		fmt.Println("Allocate Block: ",allocatedBlock)

		d.writeFS()
	}
	return allocatedBlock
}

func (d *Disk) UnallocateBlock(block int64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	if block >= d.blocks{
		panic(fmt.Sprintf("block num is bigger than cant of blocks"))
	}else if block == 0 {
		panic(fmt.Sprintf("block 0 can not be unallocated"))
	}else if d.blocks == (d.freeBlocks + 1) {
		panic(fmt.Sprintf("there are no more free blocks to unallocate"))
	}

	buffer := make([]byte,d.sizeOfBlock)
	nextBlock := make([]byte,8)
	binary.PutVarint(nextBlock, block)
	copy(buffer[d.sizeOfBlock-8:],nextBlock)

	d.volumeManager.WriteBlock(d.tailBlock,buffer)
	d.tailBlock = block

	d.freeBlocks += 1
	d.freeSpace += d.sizeOfBlock

	fmt.Println("Unallocate Block: ",block)
	d.writeFS()
}

func (d Disk) PrintInfo() {
	fmt.Println("size: ",d.size)
	fmt.Println("free space: ",d.freeSpace)
	fmt.Println("size of block: ",d.sizeOfBlock)
	fmt.Println("total blocks: ",d.blocks)
	fmt.Println("free blocks: ",d.freeBlocks)
}

func CreateDisk(diskName string, size, sizeOfBlock int64){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println(r)
        }
    }()

	newDisk := new(Disk)

	newDisk.setSize(round(size))
	newDisk.setSizeOfBlock(round(sizeOfBlock))

	if newDisk.size <= newDisk.sizeOfBlock {
		panic(fmt.Sprintf("%s", "size's Disk should be bigger than size of block"))
	}

	basicfs.CreateVolume(diskName,newDisk.size,newDisk.sizeOfBlock)
	newDisk.volumeManager = basicfs.MountVolume(diskName,sizeOfBlock)

	newDisk.freeSpace = size
	newDisk.blocks = newDisk.volumeManager.GetBlocksCant()
	newDisk.freeBlocks = newDisk.blocks

	Format(newDisk)
}

func Mount(diskName string, sizeOfBlock int64) *Disk {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r)
			fmt.Println(" or size of block should be passed: mount <diskName> <sizeOfBlock>")
		}
	}()
	mountedDisk := new(Disk)
	mountedDisk.volumeManager = basicfs.MountVolume(diskName,sizeOfBlock)

	buffer := make([]byte, sizeOfBlock)
	mountedDisk.volumeManager.ReadBlock(0,buffer)

	toStructFS(mountedDisk, buffer[:int64(unsafe.Sizeof(*mountedDisk)) - 8])

	return mountedDisk
}

func Unmount(d *Disk) {
	d.volumeManager.UnMountVolume()
	d = nil
}