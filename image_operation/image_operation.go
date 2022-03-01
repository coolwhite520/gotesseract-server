package image_operation

import (
	"errors"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

//fmt.Println(left, "\t", top, "\t", right, "\t", bottom, "\t", v.Word, "\t")
//dst := fmt.Sprintf("%s-%d.jpg", req.FilePath, k)
//fOut, _ := os.Create(dst)
//fIn, _ := os.Open(req.FilePath)
//err := image_operation.ClipImage(fIn, fOut, left, top, right, bottom, 100)
//if err != nil {
//	panic(err)
//}
//err = fOut.Close()
//if err != nil {
//	return
//}
//err = fIn.Close()
//if err != nil {
//	return
//}

func ClipImage(in io.Reader, out io.Writer, x0, y0, x1, y1, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}

	switch fm {
	case "jpg":
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{Quality: quality})
	case "png":
		switch origin.(type) {
		case *image.NRGBA:
			img := origin.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := origin.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := origin.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}
