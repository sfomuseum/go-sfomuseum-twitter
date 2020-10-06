package document

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"	
)

const JAVASCRIPT_PREFIX string = "window.YTD.tweet.part0 = "

func TrimJavaScriptPrefix(ctx context.Context, fh io.Reader) (io.ReadCloser, error) {
	return TrimPrefix(ctx, fh, JAVASCRIPT_PREFIX)
}

func TrimPrefix(ctx context.Context, fh io.Reader, prefix string) (io.ReadCloser, error) {

	offset := int64(len(prefix))
	whence := 0

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	
	_, err = br.Seek(offset, whence)

	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(br), nil
	
	// var buf bytes.Buffer
	// tee := io.TeeReader(fh, &buf)
	// return tee, nil
}
