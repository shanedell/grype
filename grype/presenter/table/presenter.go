package table

import (
	"io"
	"sort"

	"github.com/anchore/grype/grype/result"
	"github.com/anchore/syft/syft/pkg"
	"github.com/olekukonko/tablewriter"
)

// Presenter is a generic struct for holding fields needed for reporting
type Presenter struct {
	results result.Result
	catalog *pkg.Catalog
}

// NewPresenter is a *Presenter constructor
func NewPresenter(results result.Result, catalog *pkg.Catalog) *Presenter {
	return &Presenter{
		results: results,
		catalog: catalog,
	}
}

// Present creates a JSON-based reporting
func (pres *Presenter) Present(output io.Writer) error {
	rows := make([][]string, 0)

	columns := []string{"Name", "Installed", "Vulnerability", "Found-By"}
	for p := range pres.results.Enumerate() {
		row := []string{
			p.Package.Name,
			p.Package.Version,
			p.Vulnerability.ID,
			p.SearchKey,
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		_, err := io.WriteString(output, "No vulnerabilities found\n")
		return err
	}

	// sort by name, version, then type
	sort.SliceStable(rows, func(i, j int) bool {
		for col := 0; col < len(columns); col++ {
			if rows[i][0] != rows[j][0] {
				return rows[i][col] < rows[j][col]
			}
		}
		return false
	})

	table := tablewriter.NewWriter(output)

	table.SetHeader(columns)
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// these options allow for a more greppable table
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(true)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)

	// these options allow for a more human-readable (but not greppable) table
	//table.SetRowLine(true)
	//table.SetAutoMergeCells(true)
	//table.SetCenterSeparator("·") // + ┼ ╎  ┆ ┊ · •
	//table.SetColumnSeparator("│")
	//table.SetRowSeparator("─")

	table.AppendBulk(rows)
	table.Render()

	return nil
}