package rendition

import "C"

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var RENDITION_WORKERS_NUMBER = 20
var renditions WorkPool

func init() {
	imagick.Initialize()

	//Initialize the rendition workers
	renditions.workers = make(map[int]*Work)
	renditions.workItems = make([]Work, 0)
	sysDir := business.SystemDirectory()
	for i := 0; i < RENDITION_WORKERS_NUMBER; i++ {
		err := os.MkdirAll(filepath.Join(sysDir, "libreoffice", "worker_"+strconv.Itoa(i)), os.ModePerm)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatal(err.Error())
		}
		renditions.workers[i] = nil
	}
}

func Crop(image *[]byte, PositionX *int, PositionY *int, Width *int, Height *int) (result []byte, err error) {
	mw := imagick.NewMagickWand()
	var errM error
	if errM = mw.ReadImageBlob(*image); errM == nil {
		if errM = mw.SetImageFormat("png"); errM == nil {
			if errM = mw.ResizeImage(uint(*Width), uint(*Height), imagick.FILTER_LANCZOS); errM == nil {
				if errM = mw.CropImage(200, 200, (*PositionX * -1), (*PositionY * -1)); errM == nil {
					return mw.GetImageBlob(), nil
				}
			}
		}
	}
	if errM != nil {
		util.Loggify(errM)
		err = fmt.Errorf("not able to crop the suplied image")
	}
	return result, err
}

func Process(filePath *string, destination *string, forceReprocess bool) (thumbnail bool, rendition bool, processing bool, viewable bool, fileType *string, renditionType *string) {
	file, errF := os.Open(*filePath)
	if errF == nil {
		var errT error
		fileType, errT = Type(file)
		if errT != nil {
			return false, false, false, false, nil, nil
		}
		newOffset, errS := file.Seek(0, 0)
		if errS == nil && newOffset == 0 {
			switch {
			case strings.Contains(*fileType, "image/png"):
				img, errD := png.Decode(file)
				if errD == nil {
					thumbnail = Thumbnail(&img, destination, false)
				}
				viewable = true
			case strings.Contains(*fileType, "image/jpeg"):
				img, errD := jpeg.Decode(file)
				if errD == nil {
					thumbnail = Thumbnail(&img, destination, false)
				}
				viewable = true
			case strings.Contains(*fileType, "image/gif"):
				img, errD := gif.Decode(file)
				if errD == nil {
					thumbnail = Thumbnail(&img, destination, false)
				}
				viewable = true
			case strings.Contains(*fileType, "image/bmp"):
				img, errD := bmp.Decode(file)
				if errD == nil {
					renditionFile, errTF := os.Create(filepath.Join(*destination, "rendition"))
					if errTF == nil {
						errDIF := jpeg.Encode(renditionFile, img, &jpeg.Options{Quality: 100})
						if errDIF == nil {
							rendition = true
							renditionType = util.PointerString("image/jpeg")
						}
					}
					thumbnail = Thumbnail(&img, destination, false)
				}
			case strings.Contains(*fileType, "image/webp"):
				img, errD := webp.Decode(file)
				if errD == nil {
					renditionFile, errTF := os.Create(filepath.Join(*destination, "rendition"))
					if errTF == nil {
						errDIF := jpeg.Encode(renditionFile, img, &jpeg.Options{Quality: 100})
						if errDIF == nil {
							rendition = true
							renditionType = util.PointerString("image/jpeg")
						}
					}
					thumbnail = Thumbnail(&img, destination, false)
				}
			case strings.Contains(*fileType, "image/tiff"):
				img, errD := tiff.Decode(file)
				if errD == nil {
					renditionFile, errTF := os.Create(filepath.Join(*destination, "rendition"))
					if errTF == nil {
						errDIF := jpeg.Encode(renditionFile, img, &jpeg.Options{Quality: 100})
						if errDIF == nil {
							rendition = true
							renditionType = util.PointerString("image/jpeg")
						}
					}
					thumbnail = Thumbnail(&img, destination, false)
				}
			case strings.Contains(*fileType, "application/pdf"):
				processing = true
				renditions.assign(Work{
					Value:    *destination,
					WorkType: "application/pdf",
					Repeat:   false,
				}, forceReprocess)
			/*mw := imagick.NewMagickWand()
			if errR := mw.ReadImageFile(file); errR == nil {
				mw.SetIteratorIndex(0) // This being the page offset
				if errIF := mw.SetImageFormat("png"); errIF == nil {
					if errSR := mw.ThumbnailImage(190, 240); errSR == nil {
						errWI := mw.WriteImage(filepath.Join(*destination, "thumbnail"))
						if errWI == nil {
							thumbnail = true
						}
					}

				}
			}*/
			case strings.Contains(*fileType, "text/plain"):
				/*processing = true
				renditions.assign(Work{
					Value:    *destination,
					WorkType: "text/plain",
					Repeat:   false,
				}, forceReprocess)*/
			case strings.Contains(*fileType, "application/doc/xls/ppt") || strings.Contains(*fileType, "application/docx/xlsx/pptx"):
				processing = true
				renditions.assign(Work{
					Value:    *destination,
					WorkType: "application/doc/xls/ppt", //Doesn't matter what type of work, as it is non pdf and needs to be fully rerendered
					Repeat:   false,
				}, forceReprocess)
			default:
				//By default use imagick to process a rendition and thumbnail for any other imagic
				//supported formats i.e. "image/heic" will use the default case
				mw := imagick.NewMagickWand()
				var errR error
				if errR = mw.ReadImageFile(file); errR == nil {
					if errIF := mw.SetImageFormat("jpeg"); errIF == nil {
						if errWIR := mw.WriteImage(filepath.Join(*destination, "rendition")); errWIR == nil {
							rendition = true
							renditionType = util.PointerString("image/jpeg")
						}
					}
					if errIF := mw.SetImageFormat("png"); errIF == nil {
						if errIA := mw.SetImageAlpha(1.0); errIA == nil {
							if errSR := mw.ThumbnailImage(190, 190); errSR == nil {
								errWI := mw.WriteImage(filepath.Join(*destination, "thumbnail"))
								if errWI == nil {
									thumbnail = true
								}
							}
						}
					}
				} else {
					util.Loggify(errR)
				}
			}

		}
	}
	return thumbnail, rendition, processing, viewable, fileType, renditionType
}

type Work struct {
	Value    string
	WorkType string
	Repeat   bool
}
type WorkPool struct {
	mutex     sync.Mutex
	workers   map[int]*Work
	workItems []Work
}

func (wp *WorkPool) Contains(w Work) bool {
	for _, workItem := range wp.workItems {
		if workItem.Value == w.Value {
			return true
		}
	}
	return false
}

func (wp *WorkPool) ContainsInProcess(w Work) bool {
	for _, work := range wp.workers {
		if work != nil {
			if work.Value == w.Value {
				return true
			}
		}
	}
	return false
}

func (wp *WorkPool) assign(work Work, force bool) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if !force {
		//Check if a work item is not already in que or is processing
		if !wp.Contains(work) && !wp.ContainsInProcess(work) {
			for id, workInProcess := range wp.workers {
				if workInProcess == nil {
					wp.workers[id] = &work
					go Rendition(work, id)
					return
				}
			}
			wp.workItems = append(wp.workItems, work)
		}
	} else {
		if wp.Contains(work) {
			//Do nothing as the item has not started processing yet
			return
		} else if wp.ContainsInProcess(work) {
			//as it contains mark the work to be repeated
			//as the item needs to be reporcessed
			for _, workInProcess := range wp.workers {
				if workInProcess != nil {
					if workInProcess.Value == work.Value {
						workInProcess.Repeat = true
						return
					}
				}
			}
		} else {
			//Otherwise there is no work in queue and not in process
			//find a free worker or add the work to queue
			for id, workInProcess := range wp.workers {
				if workInProcess == nil {
					wp.workers[id] = &work
					go Rendition(work, id)
					return
				}
			}
			wp.workItems = append(wp.workItems, work)
		}
	}
}

func (wp *WorkPool) complete(id int, work Work) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	workInProcess := wp.workers[id]
	if workInProcess != nil && workInProcess.Repeat {
		workInProcess.Repeat = false
		wp.workers[id] = workInProcess
		go Rendition(*workInProcess, id)
		return
	}
	//Otherwise check if otherwork item exists
	//if exists pick up and process
	if len(wp.workItems) > 0 {
		var pop Work
		pop, wp.workItems = wp.workItems[0], wp.workItems[1:]
		wp.workers[id] = &pop
		go Rendition(pop, id)
	} else {
		wp.workers[id] = nil
	}
}

func GetLibreOfficeWorkerDirectory(id int) string {
	return filepath.Join(business.SystemDirectory(), "libreoffice", "worker_"+strconv.Itoa(id))
}

func AssignRenditionWork(path *string, fileType *string, forceRpeat bool) {
	renditions.assign(Work{
		Value:    *path,
		Repeat:   false,
		WorkType: *fileType,
	}, forceRpeat)
}

func Rendition(work Work, id int) {
	rendition := false
	viewable := false
	renditionType := "?"
	thumbnail := false
	if work.WorkType != "application/pdf" {
		command, sb := exec.Command("/Applications/LibreOffice.app/Contents/MacOS/soffice",
			"--headless",
			"-env:UserInstallation=file:///"+GetLibreOfficeWorkerDirectory(id),
			"--invisible",
			"--convert-to", "pdf:writer_pdf_Export", "file"), new(strings.Builder)
		command.Dir = work.Value
		command.Stdout = sb

		command.Run()
		if strings.Contains(sb.String(), "file.pdf using filter : writer_pdf_Export") {
			renditionFilePath := filepath.Join(work.Value, "rendition")
			errRen := os.Rename(filepath.Join(work.Value, "file.pdf"), renditionFilePath)
			if errRen == nil {
				rendition = true
				renditionType = "application/pdf"
			} else {
				log.Println("rendition rename not successful - " + errRen.Error())
			}
		} else {
			log.Println("document conv command not successful - " + sb.String())
		}
	} else {
		viewable = true
	}
	command, sb := exec.Command("/Applications/LibreOffice.app/Contents/MacOS/soffice",
		"--headless",
		"-env:UserInstallation=file:///"+GetLibreOfficeWorkerDirectory(id),
		"--invisible",
		"--convert-to", "png:writer_png_Export", "file"), new(strings.Builder)
	command.Dir = work.Value
	command.Stdout = sb
	command.Run()
	if strings.Contains(sb.String(), "file.png using filter : writer_png_Export") {
		thumbnailFile, errTF := os.Open(filepath.Join(work.Value, "file.png"))
		if errTF == nil {
			mw := imagick.NewMagickWand()
			if errR := mw.ReadImageFile(thumbnailFile); errR == nil {
				if errIF := mw.SetImageFormat("png"); errIF == nil {
					if errIA := mw.SetImageAlpha(1.0); errIA == nil {
						if errSR := mw.ThumbnailImage(190, 240); errSR == nil {
							errWI := mw.WriteImage(filepath.Join(work.Value, "thumbnail"))
							if errWI == nil {
								thumbnail = true
							}
						}
					}
				}
			}
		} else {
			log.Println("rendition pdf thumbnail file open not successful - " + errTF.Error())
		}
		os.Remove(filepath.Join(work.Value, "file.png"))
	} else {
		log.Println("document conv command not successful - " + sb.String())
	}

	// Update the file metadata
	metadataPath := filepath.Join(work.Value, "metadata")
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
	}
	metadata := structs.FileMetadata{}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
	}
	metadata.Thumbnail = thumbnail
	metadata.Rendition = rendition
	metadata.Processing = false
	metadata.RenditionType = util.PointerString(renditionType)
	metadata.Viewable = viewable

	metadataJson, errMJ := metadata.Json()
	if errMJ != nil || metadataJson == nil {
		util.Loggify(errMJ)
	}

	errMet = os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
	util.Loggify(errMet)

	renditions.complete(id, work)
}

func Thumbnail(img *image.Image, destination *string, rotateFlag bool) bool {
	dstImage := image.NewRGBA(image.Rect(0, 0, 190, 190))
	errTI := thumbnail(dstImage, *img)
	if errTI == nil {
		if rotateFlag {
			srcDim := dstImage.Bounds()
			roteated := image.NewRGBA(image.Rect(0, 0, srcDim.Dy(), srcDim.Dx()))
			rotate(roteated, dstImage, &RotateOptions{math.Pi / 2.0})
			dstImage = roteated
		}
		fthumbnailFile, errTF := os.Create(filepath.Join(*destination, "thumbnail"))
		if errTF == nil {
			errDIF := png.Encode(fthumbnailFile, dstImage)
			if errDIF == nil {
				return true
			}
		}
	}
	return false
}

// Rotate produces a rotated version of src, drawn onto dst.
func rotate(dst draw.Image, src image.Image, opt *RotateOptions) error {
	if dst == nil {
		return errors.New("graphics: dst is nil")
	}
	if src == nil {
		return errors.New("graphics: src is nil")
	}

	angle := 0.0
	if opt != nil {
		angle = opt.Angle
	}

	return I.Rotate(angle).TransformCenter(dst, src, Bilinear)
}

// Rotate produces a clockwise rotation transform of angle, in radians.
func (a Affine) Rotate(angle float64) Affine {
	s, c := math.Sincos(angle)
	return a.Mul(Affine{
		+c, +s, +0,
		-s, +c, +0,
		+0, +0, +1,
	})
}

// TransformCenter applies the affine transform to src and produces dst.
// Equivalent to
//
//	a.CenterFit(dst, src).Transform(dst, src, i).
func (a Affine) TransformCenter(dst draw.Image, src image.Image, i Interp) error {
	if dst == nil {
		return errors.New("graphics: dst is nil")
	}
	if src == nil {
		return errors.New("graphics: src is nil")
	}

	return a.CenterFit(dst.Bounds(), src.Bounds()).Transform(dst, src, i)
}

// CenterFit produces the affine transform, centered around the rectangles.
// It is equivalent to
//
//	I.Translate(-<center of src>).Mul(a).Translate(<center of dst>)
func (a Affine) CenterFit(dst, src image.Rectangle) Affine {
	dx := float64(dst.Min.X) + float64(dst.Dx())/2
	dy := float64(dst.Min.Y) + float64(dst.Dy())/2
	sx := float64(src.Min.X) + float64(src.Dx())/2
	sy := float64(src.Min.Y) + float64(src.Dy())/2
	return I.Translate(-sx, -sy).Mul(a).Translate(dx, dy)
}

// Translate produces a translation transform with pixel distances x and y.
func (a Affine) Translate(x, y float64) Affine {
	return a.Mul(Affine{
		1, 0, -x,
		0, 1, -y,
		0, 0, +1,
	})
}

// RotateOptions are the rotation parameters.
// Angle is the angle, in radians, to rotate the image clockwise.
type RotateOptions struct {
	Angle float64
}

func thumbnail(dst draw.Image, src image.Image) error {
	// Scale down src in the dimension that is closer to dst.
	sb := src.Bounds()
	db := dst.Bounds()
	rx := float64(sb.Dx()) / float64(db.Dx())
	ry := float64(sb.Dy()) / float64(db.Dy())
	var b image.Rectangle
	if rx < ry {
		b = image.Rect(0, 0, db.Dx(), int(float64(sb.Dy())/rx))
	} else {
		b = image.Rect(0, 0, int(float64(sb.Dx())/ry), db.Dy())
	}

	buf := image.NewRGBA(b)
	if err := Scale(buf, src); err != nil {
		return err
	}

	// Crop.
	// TODO(crawshaw): improve on center-alignment.
	var pt image.Point
	if rx < ry {
		pt.Y = (b.Dy() - db.Dy()) / 2
	} else {
		pt.X = (b.Dx() - db.Dx()) / 2
	}
	draw.Draw(dst, db, buf, pt, draw.Src)
	return nil
}

// Scale produces a scaled version of the image using bilinear interpolation.
func Scale(dst draw.Image, src image.Image) error {
	if dst == nil {
		return errors.New("graphics: dst is nil")
	}
	if src == nil {
		return errors.New("graphics: src is nil")
	}

	b := dst.Bounds()
	srcb := src.Bounds()
	if b.Empty() || srcb.Empty() {
		return nil
	}
	sx := float64(b.Dx()) / float64(srcb.Dx())
	sy := float64(b.Dy()) / float64(srcb.Dy())
	return I.Scale(sx, sy).Transform(dst, src, Bilinear)
}

func Type(file *os.File) (*string, error) {
	buf := make([]byte, 512)
	n, errB := file.Read(buf)
	if errB != nil || n < 1 {
		return nil, fmt.Errorf("unable to identify the type of the supplied file")
	}
	contentType := structs.DetectContentType(buf)
	return &contentType, nil
}

type Affine [9]float64

var I = Affine{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1,
}

// Transform applies the affine transform to src and produces dst.
func (a Affine) Transform(dst draw.Image, src image.Image, i Interp) error {
	if dst == nil {
		return errors.New("graphics: dst is nil")
	}
	if src == nil {
		return errors.New("graphics: src is nil")
	}

	// RGBA fast path.
	dstRGBA, dstOk := dst.(*image.RGBA)
	srcRGBA, srcOk := src.(*image.RGBA)
	interpRGBA, interpOk := i.(RGBA)
	if dstOk && srcOk && interpOk {
		return a.transformRGBA(dstRGBA, srcRGBA, interpRGBA)
	}

	srcb := src.Bounds()
	b := dst.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx, sy := a.pt(x, y)
			if inBounds(srcb, sx, sy) {
				dst.Set(x, y, i.Interp(src, sx, sy))
			}
		}
	}
	return nil
}

func (a Affine) transformRGBA(dst *image.RGBA, src *image.RGBA, i RGBA) error {
	srcb := src.Bounds()
	b := dst.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sx, sy := a.pt(x, y)
			if inBounds(srcb, sx, sy) {
				c := i.RGBA(src, sx, sy)
				off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
				dst.Pix[off+0] = c.R
				dst.Pix[off+1] = c.G
				dst.Pix[off+2] = c.B
				dst.Pix[off+3] = c.A
			}
		}
	}
	return nil
}

func inBounds(b image.Rectangle, x, y float64) bool {
	if x < float64(b.Min.X) || x >= float64(b.Max.X) {
		return false
	}
	if y < float64(b.Min.Y) || y >= float64(b.Max.Y) {
		return false
	}
	return true
}

func (a Affine) pt(x0, y0 int) (x1, y1 float64) {
	fx := float64(x0) + 0.5
	fy := float64(y0) + 0.5
	x1 = fx*a[0] + fy*a[1] + a[2]
	y1 = fx*a[3] + fy*a[4] + a[5]
	return x1, y1
}

// Scale produces a scaling transform of factors x and y.
func (a Affine) Scale(x, y float64) Affine {
	return a.Mul(Affine{
		1 / x, 0, 0,
		0, 1 / y, 0,
		0, 0, 1,
	})
}

// Mul returns the multiplication of two affine transform matrices.
func (a Affine) Mul(b Affine) Affine {
	return Affine{
		a[0]*b[0] + a[1]*b[3] + a[2]*b[6],
		a[0]*b[1] + a[1]*b[4] + a[2]*b[7],
		a[0]*b[2] + a[1]*b[5] + a[2]*b[8],
		a[3]*b[0] + a[4]*b[3] + a[5]*b[6],
		a[3]*b[1] + a[4]*b[4] + a[5]*b[7],
		a[3]*b[2] + a[4]*b[5] + a[5]*b[8],
		a[6]*b[0] + a[7]*b[3] + a[8]*b[6],
		a[6]*b[1] + a[7]*b[4] + a[8]*b[7],
		a[6]*b[2] + a[7]*b[5] + a[8]*b[8],
	}
}

// Interp interpolates an image's color at fractional co-ordinates.
type Interp interface {
	// Interp interpolates (x, y).
	Interp(src image.Image, x, y float64) color.Color
}

// RGBA is a fast-path interpolation implementation for image.RGBA.
// It is common for an Interp to also implement RGBA.
type RGBA interface {
	// RGBA interpolates (x, y).
	RGBA(src *image.RGBA, x, y float64) color.RGBA
}

// Gray is a fast-path interpolation implementation for image.Gray.
type Gray interface {
	// Gray interpolates (x, y).
	Gray(src *image.Gray, x, y float64) color.Gray
}

// Bilinear implements bilinear interpolation.
var Bilinear Interp = bilinear{}

type bilinear struct{}

func (i bilinear) Interp(src image.Image, x, y float64) color.Color {
	if src, ok := src.(*image.RGBA); ok {
		return i.RGBA(src, x, y)
	}
	return bilinearGeneral(src, x, y)
}

func bilinearGeneral(src image.Image, x, y float64) color.Color {
	p := findLinearSrc(src.Bounds(), x, y)
	var fr, fg, fb, fa float64
	var r, g, b, a uint32

	r, g, b, a = src.At(p.low.X, p.low.Y).RGBA()
	fr += float64(r) * p.frac00
	fg += float64(g) * p.frac00
	fb += float64(b) * p.frac00
	fa += float64(a) * p.frac00

	r, g, b, a = src.At(p.high.X, p.low.Y).RGBA()
	fr += float64(r) * p.frac01
	fg += float64(g) * p.frac01
	fb += float64(b) * p.frac01
	fa += float64(a) * p.frac01

	r, g, b, a = src.At(p.low.X, p.high.Y).RGBA()
	fr += float64(r) * p.frac10
	fg += float64(g) * p.frac10
	fb += float64(b) * p.frac10
	fa += float64(a) * p.frac10

	r, g, b, a = src.At(p.high.X, p.high.Y).RGBA()
	fr += float64(r) * p.frac11
	fg += float64(g) * p.frac11
	fb += float64(b) * p.frac11
	fa += float64(a) * p.frac11

	var c color.RGBA64
	c.R = uint16(fr + 0.5)
	c.G = uint16(fg + 0.5)
	c.B = uint16(fb + 0.5)
	c.A = uint16(fa + 0.5)
	return c
}

func (bilinear) RGBA(src *image.RGBA, x, y float64) color.RGBA {
	p := findLinearSrc(src.Bounds(), x, y)

	// Array offsets for the surrounding pixels.
	off00 := offRGBA(src, p.low.X, p.low.Y)
	off01 := offRGBA(src, p.high.X, p.low.Y)
	off10 := offRGBA(src, p.low.X, p.high.Y)
	off11 := offRGBA(src, p.high.X, p.high.Y)

	var fr, fg, fb, fa float64

	fr += float64(src.Pix[off00+0]) * p.frac00
	fg += float64(src.Pix[off00+1]) * p.frac00
	fb += float64(src.Pix[off00+2]) * p.frac00
	fa += float64(src.Pix[off00+3]) * p.frac00

	fr += float64(src.Pix[off01+0]) * p.frac01
	fg += float64(src.Pix[off01+1]) * p.frac01
	fb += float64(src.Pix[off01+2]) * p.frac01
	fa += float64(src.Pix[off01+3]) * p.frac01

	fr += float64(src.Pix[off10+0]) * p.frac10
	fg += float64(src.Pix[off10+1]) * p.frac10
	fb += float64(src.Pix[off10+2]) * p.frac10
	fa += float64(src.Pix[off10+3]) * p.frac10

	fr += float64(src.Pix[off11+0]) * p.frac11
	fg += float64(src.Pix[off11+1]) * p.frac11
	fb += float64(src.Pix[off11+2]) * p.frac11
	fa += float64(src.Pix[off11+3]) * p.frac11

	var c color.RGBA
	c.R = uint8(fr + 0.5)
	c.G = uint8(fg + 0.5)
	c.B = uint8(fb + 0.5)
	c.A = uint8(fa + 0.5)
	return c
}

func (bilinear) Gray(src *image.Gray, x, y float64) color.Gray {
	p := findLinearSrc(src.Bounds(), x, y)

	// Array offsets for the surrounding pixels.
	off00 := offGray(src, p.low.X, p.low.Y)
	off01 := offGray(src, p.high.X, p.low.Y)
	off10 := offGray(src, p.low.X, p.high.Y)
	off11 := offGray(src, p.high.X, p.high.Y)

	var fc float64
	fc += float64(src.Pix[off00]) * p.frac00
	fc += float64(src.Pix[off01]) * p.frac01
	fc += float64(src.Pix[off10]) * p.frac10
	fc += float64(src.Pix[off11]) * p.frac11

	var c color.Gray
	c.Y = uint8(fc + 0.5)
	return c
}

type bilinearSrc struct {
	// Top-left and bottom-right interpolation sources
	low, high image.Point
	// Fraction of each pixel to take. The 0 suffix indicates
	// top/left, and the 1 suffix indicates bottom/right.
	frac00, frac01, frac10, frac11 float64
}

func findLinearSrc(b image.Rectangle, sx, sy float64) bilinearSrc {
	maxX := float64(b.Max.X)
	maxY := float64(b.Max.Y)
	minX := float64(b.Min.X)
	minY := float64(b.Min.Y)
	lowX := math.Floor(sx - 0.5)
	lowY := math.Floor(sy - 0.5)
	if lowX < minX {
		lowX = minX
	}
	if lowY < minY {
		lowY = minY
	}

	highX := math.Ceil(sx - 0.5)
	highY := math.Ceil(sy - 0.5)
	if highX >= maxX {
		highX = maxX - 1
	}
	if highY >= maxY {
		highY = maxY - 1
	}

	// In the variables below, the 0 suffix indicates top/left, and the
	// 1 suffix indicates bottom/right.

	// Center of each surrounding pixel.
	x00 := lowX + 0.5
	y00 := lowY + 0.5
	x01 := highX + 0.5
	y01 := lowY + 0.5
	x10 := lowX + 0.5
	y10 := highY + 0.5
	x11 := highX + 0.5
	y11 := highY + 0.5

	p := bilinearSrc{
		low:  image.Pt(int(lowX), int(lowY)),
		high: image.Pt(int(highX), int(highY)),
	}

	// Literally, edge cases. If we are close enough to the edge of
	// the image, curtail the interpolation sources.
	if lowX == highX && lowY == highY {
		p.frac00 = 1.0
	} else if sy-minY <= 0.5 && sx-minX <= 0.5 {
		p.frac00 = 1.0
	} else if maxY-sy <= 0.5 && maxX-sx <= 0.5 {
		p.frac11 = 1.0
	} else if sy-minY <= 0.5 || lowY == highY {
		p.frac00 = x01 - sx
		p.frac01 = sx - x00
	} else if sx-minX <= 0.5 || lowX == highX {
		p.frac00 = y10 - sy
		p.frac10 = sy - y00
	} else if maxY-sy <= 0.5 {
		p.frac10 = x11 - sx
		p.frac11 = sx - x10
	} else if maxX-sx <= 0.5 {
		p.frac01 = y11 - sy
		p.frac11 = sy - y01
	} else {
		p.frac00 = (x01 - sx) * (y10 - sy)
		p.frac01 = (sx - x00) * (y11 - sy)
		p.frac10 = (x11 - sx) * (sy - y00)
		p.frac11 = (sx - x10) * (sy - y01)
	}

	return p
}

// TODO(crawshaw): When we have inlining, consider func (p *RGBA) Off(x, y) int
func offRGBA(src *image.RGBA, x, y int) int {
	return (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4
}
func offGray(src *image.Gray, x, y int) int {
	return (y-src.Rect.Min.Y)*src.Stride + (x - src.Rect.Min.X)
}
