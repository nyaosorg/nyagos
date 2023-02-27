package mains

import (
	"bytes"

	"golang.org/x/text/transform"
)

type _AtShebangFilter struct {
	normalLineFound bool
}

func (t *_AtShebangFilter) Reset() {
	t.normalLineFound = false
}

func (t *_AtShebangFilter) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	const sizeOfAtMark = 1
	for !t.normalLineFound {
		if len(src) < sizeOfAtMark {
			return nDst, nSrc, transform.ErrShortSrc
		}
		if src[0] != '@' {
			t.normalLineFound = true
			break
		}
		// make line empty
		lfPos := bytes.IndexByte(src, '\n')
		if lfPos < 0 {
			if atEOF {
				nSrc += len(src)
				return nDst, nSrc, nil
			}
			return nDst, nSrc, transform.ErrShortSrc
		}
		if len(dst) < 1 {
			return nDst, nSrc, transform.ErrShortDst
		}
		dst[0] = '\n'
		nDst++
		dst = dst[1:]
		nSrc += lfPos + 1
		src = src[lfPos+1:]
	}
	// normal line
	if len(dst) < len(src) {
		return nDst, nSrc, transform.ErrShortDst
	}
	n := copy(dst, src)
	nDst += n
	nSrc += n
	if len(dst) < len(src) {
		return nDst, nSrc, transform.ErrShortDst
	}
	return nDst, nSrc, nil
}
