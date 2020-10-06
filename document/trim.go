package document

import (
	"bytes"
	"context"
	"io"
)

const JAVASCRIPT_PREFIX string = "window.YTD.tweet.part0 = "

func TrimJavaScriptPrefix(ctx context.Context, fh io.ReadSeeker) (io.Reader, error) {
	return TrimPrefix(ctx, fh, JAVASCRIPT_PREFIX)
}

func TrimPrefix(ctx context.Context, fh io.ReadSeeker, prefix string) (io.Reader, error) {

	offset := int64(len(prefix))
	whence := 0

	_, err := fh.Seek(offset, whence)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	tee := io.TeeReader(fh, &buf)

	return tee, nil
}
