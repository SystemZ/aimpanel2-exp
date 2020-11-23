package filemanager

//Code from https://github.com/marcinwyszynski/directory_tree
import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"-"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
	Content string      `json:"content,omitempty"`
}

// Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.FileInfo, content string) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir(), content}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath   string    `json:"path"`
	Info       *FileInfo `json:"info"`
	Children   []*Node   `json:"children"`
	ParentName string    `json:"parent_name"`
	Parent     *Node     `json:"-"`
}

func (n *Node) String() string {
	j, _ := json.Marshal(n)
	return string(j)
}

// Create directory hierarchy.
// maxContentSize represents max size in kB of file which content will be showed
func NewTree(root string, limit int, maxContentSize int64) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}

	count := 0
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var content []byte
		if !info.IsDir() &&
			(info.Size()/1000) <= maxContentSize &&
			lib.StringInSlice(filepath.Ext(path), []string{".properties", ".yml"}) {

			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		}

		parents[path] = &Node{
			FullPath: strings.TrimPrefix(path, absRoot),
			Info:     fileInfoFromInterface(info, string(content)),
			Children: make([]*Node, 0),
		}

		count++
		if count == limit {
			return io.EOF
		}

		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil && err != io.EOF {
		return
	}

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			node.ParentName = parent.Info.Name
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
