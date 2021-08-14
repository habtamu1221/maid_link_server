package pkg

import "strings"

var ImageExtensions = []string{"jpeg", "png", "jpg", "gif", "btmp"}

// IsImage function checking whether the file is an image or not
func IsImage(filepath string) bool {
	extension := GetExtension(filepath)
	extension = strings.ToLower(extension)
	for _, e := range ImageExtensions {
		if e == extension {
			return true
		}
	}
	return false
}

// GetExtension function to return the extension of the File Input FileName
func GetExtension(Filename string) string {
	fileSlice := strings.Split(Filename, ".")
	if len(fileSlice) >= 1 {
		return fileSlice[len(fileSlice)-1]
	}
	return ""
}

// JPEGFileName function
func JPEGFileName(filename string) string {
	filenameSlice := strings.Split(filename, ".")
	if len(filenameSlice) > 1 {
		filenames := strings.Join(filenameSlice[:len(filenameSlice)-1], "")
		filenames += ".jpg"
		return filenames
	}
	return filename + ".jpg"
}
