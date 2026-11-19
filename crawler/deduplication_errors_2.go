package crawler

type (
	// ErrDeduplicationFailed indicates a failure in the deduplication process
	ErrDeduplicationFailed struct {
		Err        error
		Fingerprint string
	}

	// ErrExportFailed indicates a failure in the export pipeline
	ErrExportFailed struct {
		Err     error
		Fingerprint string
	}

	// ErrAlreadyProcessed indicates a URL has already been processed
	ErrAlreadyProcessed struct {
		URL string
	}

	// ErrDuplicateFingerprint indicates a duplicate fingerprint was encountered
	ErrDuplicateFingerprint struct {
		Fingerprint string
	}

	// ErrExportItemFailed indicates a failure in exporting a single item
	ErrExportItemFailed struct {
		Item    interface{}
		Err     error
		Fingerprint string
	}
)

func (e ErrDeduplicationFailed) Error() string {
	return fmt.Sprintf("failed to deduplicate URL: %s - %v", e.Fingerprint, e.Err)
}

func (e ErrExportFailed) Error() string {
	return fmt.Sprintf("failed to export URLs with fingerprint: %s - %v", e.Fingerprint, e.Err)
}

func (e ErrAlreadyProcessed) Error() string {
	return fmt.Sprintf("URL already processed: %s", e.URL)
}

func (e ErrDuplicateFingerprint) Error() string {
	return fmt.Sprintf("duplicate fingerprint encountered: %s", e.Fingerprint)
}

func (e ErrExportItemFailed) Error() string {
	return fmt.Sprintf("failed to export item: %+v - %v", e.Item, e.Err)
}
