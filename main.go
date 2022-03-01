package main

import (
	"github.com/kataras/iris/v12"
	"github.com/otiai10/gosseract/v2"
)

type ResLine struct {
	Left    int    `json:"left"`
	Top     int    `json:"top"`
	Right   int    `json:"right"`
	Bottom  int    `json:"bottom"`
	Content string `json:"content"`
}
type ResLines []ResLine

func main() {
	app := iris.Default()
	client := gosseract.NewClient()
	defer client.Close()
	app.Post("/extract", func(ctx iris.Context) {
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

	})
	app.Run(iris.Addr(":9090"))
}
