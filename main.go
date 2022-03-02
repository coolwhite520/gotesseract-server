package main

import (
	"bytes"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/otiai10/gosseract/v2"
	"os/exec"
)

type ResLine struct {
	Left    int    `json:"left"`
	Top     int    `json:"top"`
	Right   int    `json:"right"`
	Bottom  int    `json:"bottom"`
	Content string `json:"content"`
}
type ResLines []ResLine

var client *gosseract.Client

func init()  {
	client = gosseract.NewClient()
}

func main() {
	app := iris.Default()
	app.Post("/extract", handleExtract)  // ocr直接抽取文本，并带有坐标位置
	app.Post("/ocrmypdf", handleOcrPdf)  // 将不可编辑的pdf转为可编辑的
	app.Post("/ocrmyimg", handleOcrImg)  //  将图片转为可编辑的pdf
	app.Run(iris.Addr(":9090"))
}

//
func handleOcrImg(ctx iris.Context) {
	type ReqBody struct {
		Lang     string `json:"lang"`
		SrcFile string `json:"src_file"`
	}
	var req ReqBody
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  err.Error(),
		})
		return
	}
	cmdLine := fmt.Sprintf("tesseract -l %s %s %s output-prefix pdf", req.Lang, req.SrcFile, req.SrcFile)
	command := exec.Command("/bin/sh","-c",  cmdLine)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err = command.Run()
	if err != nil{
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  command.Stderr.(*bytes.Buffer).String(),
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"code": 200,
		"data": req.SrcFile + ".pdf",
		"msg":  "success",
	})
	return
}

func handleOcrPdf(ctx iris.Context) {
	type ReqBody struct {
		Lang     string `json:"lang"`
		SrcFile string `json:"src_file"`
	}
	var req ReqBody
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  err.Error(),
		})
		return
	}
	outPdfName := fmt.Sprintf("%s.ocr.pdf", req.SrcFile)
	cmdLine := fmt.Sprintf("ocrmypdf --force-ocr -l %s  %s %s", req.Lang, req.SrcFile, outPdfName)
	command := exec.Command("/bin/sh","-c",  cmdLine)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err = command.Run()
	if err != nil{
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  command.Stderr.(*bytes.Buffer).String(),
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"code": 200,
		"data": outPdfName,
		"msg":  "success",
	})
	return
}
func handleExtract(ctx iris.Context) {
	type ReqBody struct {
		Lang     string `json:"lang"`
		SrcFile string `json:"src_file"`
	}
	var req ReqBody
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  err.Error(),
		})
		return
	}
	//client.SetLanguage("eng", "deu", "jpn", "chi_sim")
	err = client.SetLanguage(req.Lang)
	if err != nil {
		err = client.SetLanguage("eng")
		if err != nil {
			ctx.JSON(map[string]interface{}{
				"code": -100,
				"msg":  err.Error(),
			})
			return
		}
	}

	err = client.SetImage(req.SrcFile)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  err.Error(),
		})
		return
	}
	boxes, err := client.GetBoundingBoxes(2)
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"code": -100,
			"msg":  err.Error(),
		})
		return
	}
	var lines ResLines
	for _, v := range boxes {
		left := v.Box.Min.X
		top := v.Box.Max.Y
		right := v.Box.Max.X
		bottom := v.Box.Min.Y
		line := ResLine{
			Left:    left,
			Top:     top,
			Right:   right,
			Bottom:  bottom,
			Content: v.Word,
		}
		lines = append(lines, line)
	}
	ctx.JSON(map[string]interface{}{
		"code": 200,
		"data": lines,
		"msg":  "success",
	})

}