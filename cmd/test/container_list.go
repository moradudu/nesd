package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/nesd/client"
	"github.com/nesd/pkg/formatter"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
	"text/template"
	"time"
)

var ContainerListCmd = &cobra.Command{
	Use:   "ps",
	Short: "list conatiner from unix socket",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		// ARGS
		containerClinet := client.NewClient()
		list, error := containerClinet.ContainerList()
		if error != nil {
			fmt.Println(error)
		}
		//buf, _ := json.Marshal(list)
		//fmt.Print(string(buf))
		FormatAndPrintContainerInfo(list, FormattingAndPrintingOptions{
			Stdout: cmd.OutOrStdout(),
			Quiet:  false,
			Size:   false,
		})
	},
}

type FormattingAndPrintingOptions struct {
	Stdout io.Writer
	// Only display container IDs.
	Quiet bool
	// Format the output using the given Go template (e.g., '{{json .}}', 'table', 'wide').
	Format string
	// Display total file sizes.
	Size bool
}

func FormatAndPrintContainerInfo(containers []types.Container, options FormattingAndPrintingOptions) error {
	w := options.Stdout
	var (
		wide bool
		tmpl *template.Template
	)
	switch options.Format {
	case "", "table":
		w = tabwriter.NewWriter(w, 4, 8, 4, ' ', 0)
		if !options.Quiet {
			printHeader := "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES"
			if options.Size {
				printHeader += "\tSIZE"
			}
			fmt.Fprintln(w, printHeader)
		}
	case "raw":
		return errors.New("unsupported format: \"raw\"")
	case "wide":
		w = tabwriter.NewWriter(w, 4, 8, 4, ' ', 0)
		if !options.Quiet {
			fmt.Fprintln(w, "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tPORTS\tNAMES\tRUNTIME\tPLATFORM\tSIZE")
			wide = true
		}
	default:
		if options.Quiet {
			return errors.New("format and quiet must not be specified together")
		}
		var err error
		tmpl, err = formatter.ParseTemplate(options.Format)
		if err != nil {
			return err
		}
	}

	for _, c := range containers {
		if tmpl != nil {
			var b bytes.Buffer
			if err := tmpl.Execute(&b, &c); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w, b.String()); err != nil {
				return err
			}
		} else if options.Quiet {
			if _, err := fmt.Fprintln(w, c.ID); err != nil {
				return err
			}
		} else {
			format := "%s\t%s\t%s\t%s\t%s\t%s\t%s"
			args := []interface{}{
				c.ID,
				c.Image,
				"",
				formatter.TimeSinceInHuman(time.Now()),
				c.Status,
				c.Ports,
				c.Names,
			}
			if wide {
				format += "\t%s\t%s\t%s\n"
				args = append(args, "", "", c.SizeRw)
			} else if options.Size {
				format += "\t%s\n"
				args = append(args, c.SizeRw)
			} else {
				format += "\n"
			}
			if _, err := fmt.Fprintf(w, format, args...); err != nil {
				return err
			}
		}

	}
	if f, ok := w.(formatter.Flusher); ok {
		return f.Flush()
	}
	return nil
}

func init() {
	rootCmd.AddCommand(ContainerListCmd)
}
