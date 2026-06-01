package input

import (
	"context"
	"io"
)

// PdfUseCase defines the PDF generation and caching operations.
type PdfUseCase interface {
	GetCurriculumPDF(ctx context.Context, language string, template string) (io.ReadCloser, string, error)
}
