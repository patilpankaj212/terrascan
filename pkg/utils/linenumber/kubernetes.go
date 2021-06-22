package linenumber

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type K8sLineNumberFinder struct {
	Elems     []*TraverseEelement
	Data      []byte
	StartLine int
}

func NewK8sLineNumberFinder(traversalPath string, data []byte, startLine int) (*K8sLineNumberFinder, error) {
	k := new(K8sLineNumberFinder)
	var err error
	k.Elems, err = getTraverseElements(traversalPath)
	if err != nil {
		return nil, err
	}
	k.Data = data
	k.StartLine = startLine

	return k, nil
}

func (k *K8sLineNumberFinder) FindLineNumber() int {
	lineNum := k.StartLine
	// get the yaml document node
	docNode := yaml.Node{}
	err := yaml.Unmarshal(k.Data, &docNode)
	if err != nil {
		fmt.Println("error while unmarshalling", err)
		return lineNum
	}

	if docNode.Kind != yaml.DocumentNode {
		fmt.Println("data is not document node")
		return lineNum
	}

	// mapping node will contain all the child nodes of the document
	if len(docNode.Content) > 0 {
		mappingNode := docNode.Content[0]
		for i := range k.Elems {
			traverseElem := k.Elems[i]
			node := getNodeValue(mappingNode, traverseElem.Name, i == (len(k.Elems)-1))
			if node == nil {
				// exact node with name not found (may not be present)
				// return last known line number
				return lineNum
			}
			lineNum = node.Line

			// if position is not nil, the found node should be a sequence node
			if traverseElem.Position != nil {
				if *traverseElem.Position <= int64(len(node.Content)) {
					node = node.Content[*traverseElem.Position]
					if node == nil {
						return lineNum
					}
					lineNum = node.Line
				}
			}
			mappingNode = node
		}
	}
	return lineNum
}

func getNodeValue(n *yaml.Node, key string, isLastTraversalElement bool) *yaml.Node {
	for i := range n.Content {
		if n.Content[i].Value == key {
			if isLastTraversalElement {
				return n.Content[i]
			}
			return n.Content[i+1]
		}
	}
	return nil
}

func getTraverseElements(traversalPath string) ([]*TraverseEelement, error) {
	traverseElems := make([]*TraverseEelement, 0)
	if len(traversalPath) == 0 {
		return nil, fmt.Errorf("traversal information not available")
	}
	pathElems := strings.Split(traversalPath, ".")
	for i := range pathElems {
		elems := strings.Split(pathElems[i], "[")
		t := TraverseEelement{
			Name: elems[0],
		}
		if len(elems) > 1 {
			x := elems[1][:len(elems[1])-1]
			j, err := strconv.ParseInt(x, 10, strconv.IntSize)
			if err != nil {
				return nil, fmt.Errorf("incorrect value for index in %s, error: %v", pathElems[i], err)
			}
			t.Position = &j
		}
		traverseElems = append(traverseElems, &t)
	}
	return traverseElems, nil
}

type TraverseEelement struct {
	Name     string
	Position *int64
}

// func main() {
// 	filePath := "/Users/pankajpatil/go/src/github.com/patilpankaj212/terrascan/test.yaml"
// 	data, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	k, err := NewK8sLineNumberFinder("spec.containers[1].securityContext", data, 1)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	num := k.FindLineNumber()
// 	fmt.Println("line number for 'spec.containers[1].securityContext' is ", num)
// 	// node := yaml.Node{}
// 	// err = yaml.Unmarshal(data, &node)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }

// 	// printNodeInfo(node)
// }

// func PrintNodeInfo(n yaml.Node) {
// 	fmt.Println(n.ShortTag())
// 	fmt.Println(n.Value)
// 	if len(n.Content) > 0 {
// 		for _, v := range n.Content {
// 			PrintNodeInfo(*v)
// 		}
// 	}
// }
