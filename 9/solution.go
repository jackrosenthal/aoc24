package main

import (
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type File struct {
	Id         int
	Blocks     int
	FreeBlocks int
}

type DiskRLE struct {
	FileId int // -1 for free blocks
	Blocks int
}

func findSlot(disk []DiskRLE, blocks int) int {
	for i, rle := range disk {
		if rle.FileId == -1 && rle.Blocks >= blocks {
			return i
		}
	}
	return -1
}

func main() {
	contents, err := os.ReadFile("input.txt")
	check(err)

	files := []File{}
	for i := 0; i < len(contents)/2; i++ {
		blocks, err := strconv.Atoi(string(contents[i*2]))
		check(err)
		freeBlocks, err := strconv.Atoi(string(contents[i*2+1]))
		files = append(files, File{i, blocks, freeBlocks})
	}

	diskBlocks := []int{}
	filesToProcess := make([]File, len(files))
	copy(filesToProcess, files)
	for len(filesToProcess) > 0 {
		for i := 0; i < filesToProcess[0].Blocks; i++ {
			diskBlocks = append(diskBlocks, filesToProcess[0].Id)
		}
		freeBlocks := filesToProcess[0].FreeBlocks
		filesToProcess = filesToProcess[1:]
		if len(filesToProcess) == 0 {
			break
		}
		for i := 0; i < freeBlocks; i++ {
			diskBlocks = append(diskBlocks, filesToProcess[len(filesToProcess)-1].Id)
			filesToProcess[len(filesToProcess)-1].Blocks--
			if filesToProcess[len(filesToProcess)-1].Blocks == 0 {
				filesToProcess = filesToProcess[:len(filesToProcess)-1]
				if len(filesToProcess) == 0 {
					break
				}
			}
		}
	}

	checksum := 0
	for i, fileId := range diskBlocks {
		checksum += i * fileId
	}
	fmt.Println(checksum)

	// Part 2
	disk := []DiskRLE{}
	for i := 0; i < len(files); i++ {
		disk = append(disk, DiskRLE{files[i].Id, files[i].Blocks})
		if files[i].FreeBlocks > 0 {
			disk = append(disk, DiskRLE{-1, files[i].FreeBlocks})
		}
	}

	for i := len(files) - 1; i >= 0; i-- {
		curIdx := 0
		for disk[curIdx].FileId != files[i].Id {
			curIdx++
		}

		// if no space, don't shift
		if findSlot(disk[:curIdx], files[i].Blocks) == -1 {
			continue
		}

		newDisk := []DiskRLE{}
		for _, rle := range disk {
			if rle.FileId == files[i].Id {
				newDisk = append(newDisk, DiskRLE{-1, rle.Blocks})
			} else {
				newDisk = append(newDisk, rle)
			}
		}
		slot := findSlot(newDisk, files[i].Blocks)
		remFreeBlocks := newDisk[slot].Blocks - files[i].Blocks
		newDiskLeft := newDisk[:slot]
		newDiskRight := make([]DiskRLE, len(newDisk)-slot-1)
		copy(newDiskRight, newDisk[slot+1:])
		disk = append(newDiskLeft, DiskRLE{files[i].Id, files[i].Blocks})
		if remFreeBlocks > 0 {
			disk = append(disk, DiskRLE{-1, remFreeBlocks})
		}
		disk = append(disk, newDiskRight...)
	}

	checksum2 := 0
	diskIdx := 0
	for _, rle := range disk {
		for i := 0; i < rle.Blocks; i++ {
			if rle.FileId != -1 {
				checksum2 += diskIdx * rle.FileId
			}
			diskIdx++
		}
	}
	fmt.Println(checksum2)
}
