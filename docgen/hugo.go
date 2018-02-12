package docgen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
	return h.writeDoc("groups", group, fm.Name, fm)
}

func (h *hugo) writeKind(group, kind string, fm *hugoKind) error {
	return h.writeDoc(group, kind, "future kind desc", fm)
}

func (h *hugo) writeDoc(category, name, content string, fm frontMatterer) error {
	logrus.WithFields(logrus.Fields{
		"category": category,
		"name":     name,
	}).Info("writing doc")
	if err := h.mkdir("content", category); err != nil {
		return errors.Wrapf(err, "create %s dir", category)
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(fm.FrontMatter()); err != nil {
		return err
	}

	buf.WriteString("\n")
	buf.WriteString(content)

	path := h.makePath("content", category, fm.Filename())
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

	rootDir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer rootDir.Close()
	names, err := rootDir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(h.makePath(name))
		if err != nil {
			return err
		}
	}
	return nil
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
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Draft    bool      `json:"draft"`
	KindName string    `json:"kind_name"`
	Versions []string  `json:"versions"`
}

type hugoKind struct {
	Title    string
	Name     string
	Versions []string
}

var _ frontMatterer = (*hugoKind)(nil)

func newKindFrontMatter(name string, versions []string) *hugoKind {
	return &hugoKind{
		Title:    name,
		Name:     name,
		Versions: versions,
	}
}

func (hk *hugoKind) FrontMatter() interface{} {
	return &kindFrontMatter{
		Title:    hk.Title,
		Date:     time.Now().UTC(),
		Draft:    false,
		KindName: hk.Name,
		Versions: hk.Versions,
	}
}

func (hk *hugoKind) Filename() string {
	return hk.Name + ".md"
}

type frontMatterer interface {
	FrontMatter() interface{}
	Filename() string
}
