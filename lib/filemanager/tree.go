package filemanager

//Code from https://github.com/marcinwyszynski/directory_tree
import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"-"`
	ModTime time.Time   `json:"-"`
	IsDir   bool        `json:"is_dir"`
	Content string      `json:"content,omitempty"`
}

// Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.FileInfo, content string) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir(), content}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}

func (n *Node) String() string {
	j, _ := json.Marshal(n)
	return string(j)
}

// Create directory hierarchy.
func NewTree(root string) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var content []byte
		if !info.IsDir() &&
			(info.Size()/1000) <= 10 &&
			lib.StringInSlice(filepath.Ext(path), []string{".properties", ".yml"}) {

			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		}

		parents[path] = &Node{
			FullPath: path,
			Info:     fileInfoFromInterface(info, string(content)),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}

	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
