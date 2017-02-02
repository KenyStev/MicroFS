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

	fmt.Printf("%v",d)

	binary.PutVarint(buffer[0:8], d.size)
	binary.PutVarint(buffer[8:16], d.freeSpace)
	binary.PutVarint(buffer[16:24], d.sizeOfBlock)
	binary.PutVarint(buffer[24:32], d.blocks)
	binary.PutVarint(buffer[32:40], d.freeBlocks)
	binary.PutVarint(buffer[40:48], d.headBlock)
	binary.PutVarint(buffer[48:56], d.tailBlock)

	fmt.Printf("%v",toStructFS(buffer))

	return buffer
}

func toStructFS(bytes []byte) *Disk{
	fs := new(Disk)
	fs.size, _ = binary.Varint(bytes[0:8])
	fs.freeSpace, _ = binary.Varint(bytes[8:16])
	fs.sizeOfBlock, _ = binary.Varint(bytes[16:24])
	fs.blocks, _ = binary.Varint(bytes[24:32])
	fs.freeBlocks,_ = binary.Varint(bytes[32:40])
	fs.headBlock,_ = binary.Varint(bytes[40:48])
	fs.tailBlock,_ = binary.Varint(bytes[48:56])

	return fs
}

func (d *Disk) writeFS() {
	fmt.Println(d.toBytesArray())
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
	initializeFreeBlocksList(newDisk)
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
	}
	return allocatedBlock
}

func CreateDisk(diskName string, size, sizeOfBlock int64){
	newDisk := new(Disk)

	newDisk.setSize(round(size))
	newDisk.setSizeOfBlock(round(sizeOfBlock))
	basicfs.CreateVolume(diskName,newDisk.size,newDisk.sizeOfBlock)
	newDisk.volumeManager = basicfs.MountVolume(diskName,sizeOfBlock)

	newDisk.freeSpace = size
	newDisk.blocks = newDisk.volumeManager.GetBlocksCant()
	newDisk.freeBlocks = newDisk.blocks

	Format(newDisk)
}