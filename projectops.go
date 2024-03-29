package coinknife

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/commandline"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/pin/lang"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GoFmt(targetProject string) {
	target := targetProject
	batName := "gofmt.bat"
	batTemplate := filepath.Join("assets", batName)
	batData := fileops.ReadFileToString(batTemplate)
	batData = strings.Replace(batData, "#TARGET_FOLDER#", target, -1)
	batFile := filepath.Join(batName)
	fileops.WriteStringToFile(batFile, batData)

	ext := &commandline.ExternalProcess{
		CommandName: batFile,
	}
	ext.Launch(true)
	ext.Wait()
}

func GoBuild(targetProject string) {
	target := targetProject
	batName := "gobuild.bat"
	batTemplate := filepath.Join("assets", batName)
	batData := fileops.ReadFileToString(batTemplate)
	batData = strings.Replace(batData, "#TARGET_FOLDER#", target, -1)
	batFile := filepath.Join(batName)
	fileops.WriteStringToFile(batFile, batData)

	ext := &commandline.ExternalProcess{
		CommandName: batFile,
	}
	ext.Launch(true)
	ext.Wait()
}

func ClearProject(target string, ignore map[string]bool) {
	pin.D("clear", target)
	files, err := ioutil.ReadDir(target)
	lang.CheckErr(err)

	for _, f := range files {
		fileName := f.Name()
		filePath := filepath.Join(target, fileName)
		if ignore[fileName] {
			pin.D("  skip", filePath)
			continue
		}
		pin.D("delete", filePath)
		err := os.RemoveAll(filePath)
		lang.CheckErr(err)
	}
	pin.D("")

}

func AppendGitIgnore(targetProject string) {
	file := filepath.Join(targetProject, ".gitignore")
	fileops.AppendStringToFile(file, "\\.idea/")
}

func ListInputProjectFiles(target string, set *Settings) []string {
	if fileops.IsFile(target) {
		lang.ReportErr("This is not a folder: %v", target)
	}

	files, err := ioutil.ReadDir(target)
	lang.CheckErr(err)
	result := []string{}
	for _, f := range files {
		fileName := f.Name()
		filePath := filepath.Join(target, fileName)

		if set.IgnoredFiles[fileName] {
			continue
		}
		if fileops.IsFolder(filePath) && !set.DoNotProcessSubfolders {
			children := ListInputProjectFiles(filePath, set)
			result = append(result, children...)
			continue
		}

		if fileops.IsFile(filePath) {
			result = append(result, filePath)
			continue
		}
	}
	result = append(result, target)
	lang.CheckErr(err)
	return result
}
