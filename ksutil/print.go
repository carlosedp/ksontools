package ksutil

import (
	"encoding/json"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Fprint prints objects to a writer in a format (yaml or json).
func Fprint(out io.Writer, objects []*unstructured.Unstructured, format string) error {
	switch format {
	case "yaml":
		return printYAML(out, objects)
	case "json":
		return printJSON(out, objects)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func printYAML(out io.Writer, objects []*unstructured.Unstructured) error {
	for _, obj := range objects {
		fmt.Fprintln(out, "---")
		buf, err := yaml.Marshal(obj.Object)
		if err != nil {
			return err
		}
		out.Write(buf)
	}

	return nil
}

func printJSON(out io.Writer, objects []*unstructured.Unstructured) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	for _, obj := range objects {
		// TODO: this is not valid framing for JSON
		if len(objects) > 1 {
			fmt.Fprintln(out, "---")
		}
		if err := enc.Encode(obj); err != nil {
			return err
		}
	}

	return nil
}
