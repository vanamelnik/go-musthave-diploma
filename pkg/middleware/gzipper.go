package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// gzipReadCloser replaces body and decompress it while reading.
type gzipReadCloser struct {
	gr   io.ReadCloser
	body io.ReadCloser
}

func (g gzipReadCloser) Read(p []byte) (int, error) {
	return g.gr.Read(p)
}

func (g gzipReadCloser) Close() error {
	if err := g.gr.Close(); err != nil {
		return err
	}
	return g.body.Close()
}

// gzipWriter replaces http.ResponseWriter and compresses incoming data.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (gw gzipWriter) Write(data []byte) (int, error) {
	return gw.Writer.Write(data)
}

// GzipMdlw decompresses request body if it's compressed. It also checks whether the frontend accepts gzip encoding,
// and, if so, compresses the response.
func GzipMdlw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Override request body if needed.
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("gzipHandle: %v", err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)

				return
			}

			r.Body = gzipReadCloser{
				gr:   gr,
				body: r.Body,
			}
		}

		// Overrride response writer, if needed.
		respWriter := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				log.Printf("gzipHandle: %v", err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			defer gw.Close()

			respWriter = gzipWriter{ResponseWriter: w, Writer: gw}
			respWriter.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(respWriter, r)
	})
}
