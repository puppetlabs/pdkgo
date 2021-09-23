package main

import (
	"net/http"

	"github.com/puppetlabs/pdkgo/cmd/pct/build"
	"github.com/puppetlabs/pdkgo/cmd/pct/completion"
	"github.com/puppetlabs/pdkgo/cmd/pct/install"
	"github.com/puppetlabs/pdkgo/cmd/pct/new"
	"github.com/puppetlabs/pdkgo/cmd/pct/root"
	appver "github.com/puppetlabs/pdkgo/cmd/pct/version"
	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var rootCmd = root.CreateRootCommand()

	var verCmd = appver.CreateVersionCommand(version, date, commit)
	v := appver.Format(version, date, commit)
	rootCmd.Version = v
	rootCmd.SetVersionTemplate(v)
	rootCmd.AddCommand(verCmd)

	rootCmd.AddCommand(completion.CreateCompletionCommand())

	// afero setup
	fs := afero.NewOsFs()
	afs := afero.Afero{Fs: fs}
	iofs := afero.IOFS{Fs: fs}

	// build
	rootCmd.AddCommand(build.CreateCommand())

	// install
	installCmd := install.InstallCommand{
		PctInstaller: &pct.PctInstaller{
			Tar:        &tar.Tar{AFS: &afs},
			Gunzip:     &gzip.Gunzip{AFS: &afs},
			AFS:        &afs,
			IOFS:       &iofs,
			HTTPClient: &http.Client{},
		},
	}
	rootCmd.AddCommand(installCmd.CreateCommand())

	// new
	rootCmd.AddCommand(new.CreateCommand())

	cobra.OnInitialize(root.InitLogger, root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}