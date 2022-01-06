package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/otiai10/gosseract/v2"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)



func main() {
	app := iris.Default()
	os.MkdirAll("./uploads", 0666)
	app.Post("/upload", func(ctx iris.Context) {
		client := gosseract.NewClient()
		// lang 简称
		lang := ctx.FormValue("lang")
		//client.SetLanguage("eng", "deu", "jpn", "chi_sim")
		err := client.SetLanguage(lang)
		if err != nil {
			client.SetLanguage("eng")
		}
		defer client.Close()
		file, info, err := ctx.FormFile("image")
		if err != nil {
			ctx.JSON(map[string]interface{}{
				"code": -100,
				"msg": err.Error(),
				"content": "",
			})
			return
		}
		defer file.Close()
		newFilename := strconv.Itoa(int(time.Now().UnixNano())) + path.Ext(info.Filename)
		filePathName := fmt.Sprintf("./uploads/%s", newFilename)
		out, err := os.OpenFile(filePathName, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			ctx.JSON(map[string]interface{}{
				"code": -100,
				"msg": err.Error(),
				"content": "",
			})
			return
		}
		defer out.Close()
		io.Copy(out, file)
		client.SetImage(filePathName)
		text, err := client.Text()
		os.Remove(filePathName)
		if err != nil {
			ctx.JSON(map[string]interface{}{
				"code": -100,
				"msg": err.Error(),
				"content": "",
			})
			return
		}
		ctx.JSON(map[string]interface{}{
			"code": 200,
			"content": text,
		})

	})
	app.Run(iris.Addr(":9090"))
}
