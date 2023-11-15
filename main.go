package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

const FILE_NAME = "tree.gob"

type TreeNode struct {
    Key   string
    Value string
    Left  *TreeNode
    Right *TreeNode
}

func (t *TreeNode) Insert(key string, value string) {
    if t == nil {
        t = &TreeNode{Key: key, Value: value}
    } else if key < t.Key {
        if t.Left == nil {
            t.Left = &TreeNode{Key: key, Value: value}
        } else {
            t.Left.Insert(key, value)
        }
    } else {
        if t.Right == nil {
            t.Right = &TreeNode{Key: key, Value: value}
        } else {
            t.Right.Insert(key, value)
        }
    }
}

func (t *TreeNode) Search(key string) *TreeNode {
    if t == nil || t.Key == key {
        return t
    }
    if key < t.Key {
        return t.Left.Search(key)
    }
    return t.Right.Search(key)
}

func (t *TreeNode) Delete(key string) *TreeNode {
    if t == nil {
        return nil
    }
    if key < t.Key {
        t.Left = t.Left.Delete(key)
    } else if key > t.Key {
        t.Right = t.Right.Delete(key)
    } else {
        if t.Left == nil {
            return t.Right
        } else if t.Right == nil {
            return t.Left
        }
        minNode := t.Right.MinNode()
        t.Key = minNode.Key
        t.Value = minNode.Value
        t.Right = t.Right.Delete(t.Key)
    }
    return t
}

func (t *TreeNode) MinNode() *TreeNode {
    current := t
    for current.Left != nil {
        current = current.Left
    }
    return current
}

func (t *TreeNode) Update(key string, value string) {
    node := t.Search(key)
    if node != nil {
        node.Value = value
    }
}

func inorderTraversal(root *TreeNode) {
    if root != nil {
        inorderTraversal(root.Left)
        fmt.Println(root.Key)
        inorderTraversal(root.Right)
    }
}

func loadTree(filename string) (*TreeNode, error) {
    file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, nil
	}

    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := gob.NewDecoder(file)
    var tree TreeNode
    err = decoder.Decode(&tree)
    if err != nil {
        return nil, err
    }

    return &tree, nil
}

func saveTree(tree *TreeNode, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := gob.NewEncoder(file)
    err = encoder.Encode(tree)
    if err != nil {
        return err
    }

    return nil
}

func main() {

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Invalid number of arguments")
		return
	}

	tree, err := loadTree(FILE_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}
	if tree == nil {
		tree = &TreeNode{}
	}

	operation := args[0]
	key := args[1]
	var value string
	if len(args) > 2 {
		value = args[2]
	}

	switch operation {
	case "list": {
		inorderTraversal(tree)
	}
	case "insert":
		tree.Insert(key, value)
	case "query":
		res := tree.Search(key)
		if res != nil {
			fmt.Println(res.Value)
			break
		}
		fmt.Println("Cannot find key: ", key)
	case "delete":
		res := tree.Search(key)
		if res != nil {
			tree.Delete(key)
			break
		}
		fmt.Println("Cannot find key: ", key)
	case "update":
		res := tree.Search(key)
		if res != nil {
			tree.Update(key, value)
			break
		}
		fmt.Println("Cannot find key: ", key)
	default:
		fmt.Println("Invalid operation, we only support query/insert/update/delete [key] [value]")
	}
	saveTree(tree, FILE_NAME)
}