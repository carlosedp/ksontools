package docgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type hugo struct {
	root string
}

func newHugo(root string) (*hugo, error) {

	h := &hugo{
		root: root,
	}

	if err := h.cleanContents(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *hugo) makePath(path ...string) string {
	return filepath.Join(append([]string{h.root}, path...)...)
}

func (h *hugo) mkdir(path ...string) error {
	dirName := h.makePath(path...)
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return errors.Wrapf(err, "create directory %s", dirName)
	}

	return nil
}

func (h *hugo) writeGroup(group string, fm *hugoGroup) error {
	return h.writeDoc([]string{"groups"}, group, fm.Name, fm)
}

func (h *hugo) writeKind(group, kind string, fm *hugoKind) error {
	return h.writeDoc([]string{group}, kind, "future kind desc", fm)
}

func (h *hugo) writeVersionedKind(group, version, kind string) error {
	category := []string{group, version}
	fm := newHugoVersionedKind(version, kind)
	content := fmt.Sprintf("%s/%s/%s", group, version, kind)
	return h.writeDoc(category, kind, content, fm)
}

func (h *hugo) writeDoc(category []string, name, content string, fm frontMatterer) error {
	logrus.WithFields(logrus.Fields{
		"category": strings.Join(category, "/"),
		"name":     name,
	}).Info("writing doc")

	parentPath := append([]string{"content"}, category...)
	if err := h.mkdir(parentPath...); err != nil {
		return errors.Wrapf(err, "create %s dir", category)
	}

	var buf bytes.Buffer

	b, err := json.MarshalIndent(fm.FrontMatter(), "", "  ")
	if err != nil {
		return err
	}

	if _, err := buf.Write(b); err != nil {
		return err
	}

	buf.WriteString("\n")
	buf.WriteString(content)

	path := h.makePath(append(parentPath, fm.Filename())...)
	// path := h.makePath("content", category, fm.Filename())
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}

func (h *hugo) cleanContents() error {
	path := h.makePath("content")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrap(err, "check content path")
	}

	return os.RemoveAll(path)
}

type versionedKindFrontMatter struct {
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	Draft  bool      `json:"draft"`
	Layout string    `json:"layout"`
	Type   string    `json:"type"`
}

type hugoVersionedKind struct {
	version string
	kind    string
}

var _ frontMatterer = (*hugoVersionedKind)(nil)

func newHugoVersionedKind(version, kind string) *hugoVersionedKind {
	return &hugoVersionedKind{
		version: version,
		kind:    kind,
	}
}

func (hvk *hugoVersionedKind) FrontMatter() interface{} {
	return &versionedKindFrontMatter{
		Title:  hvk.kind,
		Date:   time.Now().UTC(),
		Draft:  false,
		Layout: "kind",
		Type:   "kind",
	}
}

func (hvk *hugoVersionedKind) Filename() string {
	return hvk.kind + ".md"
}

func newGroupFrontMatter(name string) *hugoGroup {
	return &hugoGroup{
		Title: name,
		Name:  name,
	}
}

type groupFrontMatter struct {
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Draft     bool      `json:"draft"`
	GroupName string    `json:"group_name"`
}

type hugoGroup struct {
	Title string
	Name  string
}

var _ frontMatterer = (*hugoGroup)(nil)

func (hg *hugoGroup) FrontMatter() interface{} {
	return &groupFrontMatter{
		Title:     hg.Title,
		Date:      time.Now().UTC(),
		Draft:     false,
		GroupName: hg.Name,
	}
}

func (hg *hugoGroup) Filename() string {
	return hg.Name + ".md"
}

type kindFrontMatter struct {
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Draft       bool      `json:"draft"`
	KindName    string    `json:"kind_name"`
	Versions    []string  `json:"versions"`
	ParentGroup string    `json:"parent_group"`
}

type hugoKind struct {
	Title    string
	group    string
	Name     string
	Versions []string
}

var _ frontMatterer = (*hugoKind)(nil)

func newKindFrontMatter(group, name string, versions []string) *hugoKind {
	return &hugoKind{
		Title:    name,
		group:    group,
		Name:     name,
		Versions: versions,
	}
}

func (hk *hugoKind) FrontMatter() interface{} {
	return &kindFrontMatter{
		Title:       hk.Title,
		Date:        time.Now().UTC(),
		Draft:       false,
		KindName:    hk.Name,
		Versions:    hk.Versions,
		ParentGroup: hk.group,
	}
}

func (hk *hugoKind) Filename() string {
	return hk.Name + ".md"
}

type frontMatterer interface {
	FrontMatter() interface{}
	Filename() string
}
