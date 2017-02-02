package basicfs

import (
	"os"
	"log"
)

const rootpath = "disks/"

type Volume struct{
	mountedDisk *os.File
	sizeOfBlock int64
}

func (fsl *Volume) setDisk(f *os.File) {
	fsl.mountedDisk = f
}

func (fsl *Volume) setSizeOfBlock(s int64) {
	fsl.sizeOfBlock = s
}

func openfile(filepath string) *os.File {
	file,err := os.OpenFile(rootpath + filepath, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil{
		log.Println("no abrio archivo "+filepath)
		return nil
	}
	return file
}

func (fsl *Volume) GetBlocksCant() int64 {
	fileinfo,_ := fsl.mountedDisk.Stat()
	return fileinfo.Size()/fsl.sizeOfBlock
}

// Public Funcs

func CreateVolume(volumeName string, size, sizeOfBlock int64) {
	disk := openfile(volumeName)
	defer disk.Close()

	buffer := make([]byte,sizeOfBlock)
	for i := 0; i < int(size/sizeOfBlock); i++ {
		disk.Write(buffer)
	}
}

func MountVolume(volumeName string, sizeOfBlock int64) *Volume{
	newMountedVolum := new(Volume)
	newMountedVolum.setSizeOfBlock(sizeOfBlock)
	newMountedVolum.setDisk(openfile(volumeName))

	return newMountedVolum
}

// Public Volume funcs 

func (fsl *Volume) UnMountVolume() {
	fsl.mountedDisk.Close()
}

func (fsl *Volume) ReadBlock(numBlock int64, buffer []byte) {
	fsl.mountedDisk.Seek(numBlock*fsl.sizeOfBlock,0)
	fsl.mountedDisk.Read(buffer)
}

func (fsl *Volume) WriteBlock(numBlock int64, buffer []byte) {
	fsl.mountedDisk.Seek(numBlock*fsl.sizeOfBlock,0)
	fsl.mountedDisk.Write(buffer)
}