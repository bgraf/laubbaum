package model

type Node struct {
	Text     string
	Children []*Node
	Parent   *Node
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) SetText(newText string) {
	n.Text = newText
}

func (n *Node) InnerSize() (w, h int) {
	return len(n.Text), 1
}

func (n *Node) PropagateParent(parent *Node) {
	n.Parent = parent

	for _, c := range n.Children {
		c.PropagateParent(n)
	}
}

func (n *Node) ChildIndex(child *Node) (int, bool) {
	for i, c := range n.Children {
		if c == child {
			return i, true
		}
	}

	return 0, false
}

func (n *Node) AppendChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) InsertChild(child *Node, position int) {
	if position >= len(n.Children) {
		n.AppendChild(child)

		return
	}

	n.Children = append(n.Children, nil)
	copy(n.Children[position+1:], n.Children[position:])

	n.Children[position] = child
	child.Parent = n
}

func (n *Node) PreviousSibling() *Node {
	if n.Parent == nil {
		return nil
	}

	index, _ := n.Parent.ChildIndex(n)
	if index > 0 {
		return n.Parent.Children[index-1]
	}

	return nil
}

func (n *Node) NextSibling() *Node {
	if n.Parent == nil {
		return nil
	}

	index, _ := n.Parent.ChildIndex(n)
	if index+1 < len(n.Parent.Children) {
		return n.Parent.Children[index+1]
	}

	return nil
}
