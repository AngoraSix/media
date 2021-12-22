package media

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"os"
	"strings"

	images "angorasix.com/media/gen/images"
)

// ImagesUploadDecoderFunc implements the multipart decoder for service
// "images" endpoint "upload". The decoder must populate the argument p after
// encoding.
func ImagesUploadDecoderFunc(mr *multipart.Reader, p **images.ImageUploadPayload) error {
	res := images.ImageUploadPayload{}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		_, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
		if err != nil {
			// can't process this entry, it probably isn't an image
			continue
		}

		disposition, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		// the disposition can be, for example 'image/jpeg' or 'video/mp4'
		// We want to support only image files!
		if err != nil || !strings.HasPrefix(disposition, "image/") {
			// can't process this entry, it probably isn't an image
			continue
		}

		if params["name"] == "file" {
			bytes, err := ioutil.ReadAll(p)
			if err != nil {
				// can't process this entry, for some reason
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			filename := params["filename"]
			imageUpload := images.ImageUpload{
				Type:  &disposition,
				Bytes: bytes,
				Name:  &filename,
			}
			res.Files = append(res.Files, &imageUpload)
		}
	}
	*p = &res
	return nil
}

// ImagesUploadEncoderFunc implements the multipart encoder for service
// "images" endpoint "upload".
func ImagesUploadEncoderFunc(mw *multipart.Writer, p *images.ImageUploadPayload) error {
	// Add multipart request encoder logic here
	return nil
}
