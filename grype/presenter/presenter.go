package presenter

import (
	"io"

	"github.com/anchore/grype/grype/presenter/template"

	"github.com/anchore/grype/grype/match"
	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/cyclonedx"
	"github.com/anchore/grype/grype/presenter/json"
	"github.com/anchore/grype/grype/presenter/table"
	"github.com/anchore/grype/grype/vulnerability"
)

// Presenter is the main interface other Presenters need to implement
type Presenter interface {
	Present(io.Writer) error
}

// GetPresenter retrieves a Presenter that matches a CLI option
func GetPresenter(presenterConfig Config, matches match.Matches, packages []pkg.Package, context pkg.Context,
	metadataProvider vulnerability.MetadataProvider, appConfig interface{}) Presenter {
	switch presenterConfig.format {
	case jsonFormat:
		return json.NewPresenter(matches, packages, context, metadataProvider, appConfig)
	case tableFormat:
		return table.NewPresenter(matches, packages, metadataProvider)
	case cycloneDXFormat:
		return cyclonedx.NewPresenter(matches, packages, context.Source, metadataProvider)
	case templateFormat:
		return template.NewPresenter(matches, packages, context, metadataProvider, appConfig, presenterConfig.templateFilePath)
	default:
		return nil
	}
}
