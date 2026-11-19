package crawler

import (
	"fmt"
	"sync"

	"github.com/olekukonko/tablewriter"
)

// URLFrontierFormatter formats the URL frontier data in a human-readable table.
type URLFrontierFormatter struct {
	table *tablewriter.Table
	mx    sync.Mutex
}

// NewURLFrontierFormatter returns a new instance of URLFrontierFormatter.
func NewURLFrontierFormatter() *URLFrontierFormatter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeaderAlign(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	return &URLFrontierFormatter{table: table}
}

// Format updates the table with the latest URL frontier data and prints it.
func (f *URLFrontierFormatter) Format(urls []string) {
	f.mx.Lock()
	defer f.mx.Unlock()

	f.table.SetHeader([]string{"Index", "URL"})
	for i, url := range urls {
		f.table.Append([]string{fmt.Sprintf("%d", i), url})
	}
	f.table.Render()
}
