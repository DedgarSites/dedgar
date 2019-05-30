package tree

import (
	"fmt"
	"path"
	"strings"
)

var (
	startPath  = "/"
	RootFolder = newFolder(startPath)
)

type File struct {
	//Id   string
	Name string
}

type Folder struct {
	Name    string
	Files   []File
	Folders map[string]*Folder
}

func newFolder(name string) *Folder {
	return &Folder{name, []File{}, make(map[string]*Folder)}
}

func (f *Folder) getFolder(name string) *Folder {
	if nextF, ok := f.Folders[name]; ok {
		return nextF
	} else if f.Name == name {
		return f
	} else {
		return &Folder{}
	}
}

func (f *Folder) existFolder(name string) bool {
	for _, v := range f.Folders {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (f *Folder) addFolder(folderName string) {
	if !f.existFolder(folderName) {
		f.Folders[folderName] = newFolder(folderName)
	}
}

func (f *Folder) addFile(fileName string) {
	f.Files = append(f.Files, File{fileName})
}

func (f *Folder) getList() (result []map[string]interface{}) {
	for _, v := range f.Folders {
		result = append(result, map[string]interface{}{
			"name": v.Name,
			"type": "folder",
		})
	}

	for _, v := range f.Files {
		result = append(result, map[string]interface{}{
			//"id":   v.Id,
			"name": v.Name,
			"type": "file",
		})
	}
	return
}

func isFile(str string) bool {
	if path.Ext(str) != "" {
		return true
	}
	return false
}

func DeleteEmptyElements(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

type IS map[string]string

func init() {
	arrayPaths := []interface{}{
		IS{
			"id":       "1",
			"filePath": "all/",
		},
		IS{
			"id":       "1",
			"filePath": "all/peelz.here",
		},
		IS{
			"id":       "2",
			"filePath": "all/test1/",
		},
		IS{
			"id":       "3",
			"filePath": "all/test1/Nene_noises_for_1_32_minutes.mp4",
		},
		IS{
			"id":       "3",
			"filePath": "all/test1/neptune_all_the_meme.jpg",
		},
		IS{
			"id":       "3",
			"filePath": "all/test2/",
		},
		IS{
			"id":       "3",
			"filePath": "all/test2/america_chan_seijouki.png",
		},
		IS{
			"id":       "3",
			"filePath": "all/test2/bongo_cat_levan_polka_miku.mp4",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/inside_test3.jpg",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/test4/",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/test4/second_level.jpg",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/test4/another_s2.mp3",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/test4/test5/",
		},
		IS{
			"id":       "3",
			"filePath": "all/test3/test4/test5/test6/",
		},
	}

	breadcrumb := "all/"

	for _, path := range arrayPaths {
		filePath := path.(IS)["filePath"]
		//fileId := path.(IS)["id"]
		splitPath := DeleteEmptyElements(strings.Split(filePath, "/"))
		tmpFolder := RootFolder
		for _, item := range splitPath {
			if isFile(item) {
				tmpFolder.addFile(item) //, fileId)
			} else {
				if item != startPath {
					tmpFolder.addFolder(item)
				}
				tmpFolder = tmpFolder.getFolder(item)
			}
		}
	}

	currentFolder := RootFolder.getFolder("/")
	breadcrumbElements := DeleteEmptyElements(strings.Split(breadcrumb, "/"))
	for i, v := range breadcrumbElements {
		if currentFolder.existFolder(v) {
			currentFolder = currentFolder.getFolder(v)
			if i == len(breadcrumbElements)-1 {
				break
			}
		} else {
			currentFolder = currentFolder.getFolder(v)
		}
	}
	//fmt.Println(currentFolder.getList())
	//printDir(RootFolder)
}

func printDir(RootFolder *Folder) {
	fmt.Println("In folder: ", RootFolder.Name)
	fmt.Println(RootFolder.Files)
	/*if RootFolder.Folders != nil {
		for _, fol := range RootFolder.Folders {
			fmt.Println(fol.Name)
		}
	}*/

	if len(RootFolder.Folders) > 0 {
		for _, folder := range RootFolder.Folders {
			fmt.Printf("found nested folder: %+v", folder)
			printDir(folder)
		}
	}
	//fmt.Println(RootFolder.Name)
}

func FindObj(rootFolder *Folder, findItem string, fChan chan *Folder, eChan chan bool) {
	fmt.Println("In folder: ", rootFolder.Name)
	select {
	case <-eChan:
		return
	case folderChan := <-fChan:
		fmt.Println("returning", folderChan.Name)
		eChan <- true
		return
	default:
	}
	if rootFolder.Name == findItem {
		fmt.Println("this is already what we're looking for")
		fChan <- rootFolder
		return
	}
	for _, folder := range rootFolder.Folders {
		if folder.Name == findItem {
			fmt.Printf("one of the folders in %v is what we're looking for\n", rootFolder.Name)
			fChan <- folder
			return
		}
		go FindObj(folder, findItem, fChan, eChan)
	}
}

func FindNode(rootFolder *Folder, findItem string) *Folder {
	found := newFolder("")

	if rootFolder.Name == findItem {
		return rootFolder
	}

	for _, folder := range rootFolder.Folders {
		if folder.Name == findItem {
			return folder
		}
		found = FindNode(folder, findItem)
	}

	return found
}
