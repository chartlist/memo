// Package archive /**
package archive

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/urfave/cli"
	"os"
)

var ArchiveCommand = cli.Command{
	Name:   "archive",
	Action: doArchive,
	Usage:  "do archive",
}

func doArchive(ctx *cli.Context) {

	dir := "../www/archives/text" // 文件夹路径

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	type message struct {
		Id   string
		Body string
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("bitcoin.archive.db", mapping)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() { // 如果是文件
			fileName := file.Name()
			fileContent, err := os.ReadFile(dir + "/" + fileName)
			if err != nil {
				panic(err)
			}
			msg := message{
				Id:   fileName,
				Body: string(fileContent),
			}
			err = index.Index(msg.Id, msg)
			if err != nil {
				fmt.Println("index.Index err:", err.Error())
			} else {
				fmt.Println("index ", fileName, " success")
			}

			//bytes += len(fileContent)
			//fmt.Println(fileName+":\n", len(fileContent))
		}
	}
	fmt.Println("all Finished.....")
	//fmt.Println("bytes:", bytes)
}
