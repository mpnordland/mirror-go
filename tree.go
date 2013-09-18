package main

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/xrect"
)

//Node is an interface that allows
//Leaf and Branch to be used interchangably
//it embeds the Signalable interface
type Node interface {
	Signalable
	AddChild(n Node)
	RemoveChild(n Node)
	Parent() Node
	ParentSet(n Node)
	SplitType() string
	SplitTypeSet(sType string)
	Layout(rect *xrect.XRect)
	Children() NodeSlice
	VisitLeaves(f func(n Node) bool) bool
	VisitLeavesReverse(f func(n Node) bool) bool
	Print(indent string)
}

//NodeSlice is just a rename of []Node
//originally it was going to be extened
//to match sort.Interface to take advantage
//of sort.Reverse(), but then I realized
//that it did not reverse the items in place
type NodeSlice []Node

//sliceRemove is a convinience function for
//removing an item from a NodeSlice
func sliceRemove(ns NodeSlice, i int) NodeSlice {
	return append(ns[:i], ns[i+1:]...)
}

//Leaf is a struct that represents the window
//parts of the tree.
type Leaf struct {
	*Dispatcher
	parent Node
	win    *Window
}

//NewLeaf creates a new Leaf, attaching
//it's removeSelf method to the "window::unmapped"
//hook of the window it holds
func NewLeaf(win *Window) *Leaf {
	l := &Leaf{NewDispatcher(), nil, win}
	win.AttachToHook("window::unmapped", l.removeSelf)
	return l
}

//AddChild makes Tree.Insert easier
//trying to add a child makes the logical
//assumtion that you really want to split
//the space taken by that Leaf with another
//Node. It does so by creating a new parent
//Branch and setting itself and the new Node
//as its children
func (l *Leaf) AddChild(n Node) {
	lParent := l.parent
	l.parent.RemoveChild(l)
	b := NewBranch()
	if lParent.SplitType() == "horizontal" {
		b.SplitTypeSet("vertical")
	} else {
		b.SplitTypeSet("horizontal")
	}
	b.AddChild(l)
	b.AddChild(n)
	lParent.AddChild(b)
}

//RemoveChild does nothing as a Leaf
//has no children.
func (l *Leaf) RemoveChild(n Node) {}

//SplitType returns an empty string
//it's here to fulfil the Node interface{}
func (l *Leaf) SplitType() string {
	return ""
}

//SplitTypeSet does nothing, it's here
//to fulfil the Node interface
func (l *Leaf) SplitTypeSet(sType string) {}

//Layout passes the XRect it recieves to
//the window it holds
func (l *Leaf) Layout(rect *xrect.XRect) {
	l.win.GeometrySet(rect)
}

//Parent returns the parent of this leaf
func (l *Leaf) Parent() Node {
	return l.parent
}

//ParentSet sets the parent of this Leaf
func (l *Leaf) ParentSet(n Node) {
	l.parent = n
}

//VisitLeaves calls the passed function on itself
//and returns the result. See Branch.VisitLeaves
//for more information
func (l *Leaf) VisitLeaves(f func(n Node) bool) bool {
	return f(l)
}

//VisitLeavesReverse calls the passed function on itself
//and returns the result. See Branch.VisitLeavesReverse
//for more information
func (l *Leaf) VisitLeavesReverse(f func(n Node) bool) bool {
	return f(l)
}

//Children returns an empty NodeSlice
//It's here to fulfil the Node interface
func (l *Leaf) Children() NodeSlice {
	return make(NodeSlice, 0)
}

//removeSelf is used to ask the tree to properly
//remove this Leaf
func (l *Leaf) removeSelf(args ...interface{}) {
	l.PullHook("leaf::remove", l)
}

//Print prints this Leaf along with its window
func (l *Leaf) Print(indent string) {
	fmt.Printf("%sLeaf -- %d\n", indent, l.win.Id)
}

//Branch is a struct used for representing the purely
//stuctural parts of the tree, i.e. anything that
//simply describes arrangement, and does not directly
//hold a window
type Branch struct {
	*Dispatcher
	parent    Node
	children  NodeSlice
	splitType string
}

//NewBranch creates a new Branch
func NewBranch() *Branch {
	return &Branch{NewDispatcher(), nil, make(NodeSlice, 0), "horizontal"}
}

//AddChild adds the node passed to it as
//a child of this Branch
func (b *Branch) AddChild(n Node) {
	n.ParentSet(b)
	b.children = append(b.children, n)
}

//RemoveChild removes the child passed to it
func (b *Branch) RemoveChild(n Node) {
	for i, nn := range b.children {
		if nn == n {
			b.children = sliceRemove(b.children, i)
			b.children = append(b.children, n.Children()...)
			break
		}
	}
}

//SplitType returns the split type of this Branch
func (b *Branch) SplitType() string {
	return b.splitType
}

//SplitTypeSet sets the split type of this Branch
func (b *Branch) SplitTypeSet(sType string) {
	b.splitType = sType
}

//Layout splits the XRect passed to it
//according to its split type and number of children
//and calls the Layout function of each of its children
//with one of the sections
func (b *Branch) Layout(rect *xrect.XRect) {
	if len(b.children) == 0 {
		return
	}
	if b.splitType == "horizontal" {
		width := rect.Width() / len(b.children)
		x := rect.X()
		for _, n := range b.children {
			n.Layout(xrect.New(x, rect.Y(), width, rect.Height()))
			x += width
		}
	} else {
		height := rect.Height() / len(b.children)
		y := rect.Y()
		for _, n := range b.children {
			n.Layout(xrect.New(rect.X(), y, rect.Width(), height))
			y += height
		}
	}
}

//Parent returns the parent of this Branch
func (b *Branch) Parent() Node {
	return b.parent
}

//ParentSet sets the parent of this Branch
func (b *Branch) ParentSet(n Node) {
	b.parent = n
}

//VisitLeaves is used for forward walking the tree leaves.
//If the function passed to it returns false, walking stops.
//If it returns true walking continues.
func (b *Branch) VisitLeaves(f func(n Node) bool) bool {
	for _, c := range b.children {
		if !c.VisitLeaves(f) {
			return false
		}
	}
	return true
}

//VisitLeavesReverse is used for reverse walking the tree
//leaves to find a particular one. If the function passed to
//it returns false, walking stops. If it returns true walking
//continues
func (b *Branch) VisitLeavesReverse(f func(n Node) bool) bool {
	for i := len(b.children) - 1; i >= 0; i-- {
		c := b.children[i]
		if !c.VisitLeavesReverse(f) {
			return false
		}
	}
	return true
}

//Children returns the children of this
//Branch.
func (b *Branch) Children() NodeSlice {
	return b.children
}

//Print prints this Branch with an indent
//and then calls each of its children's
//Print method
func (b *Branch) Print(indent string) {
	fmt.Printf("%sBranch\n", indent)
	indent += " "
	for _, c := range b.children {
		c.Print(indent)
	}
}

//Tree is a struct that represents
//a tree for arranging windows.
//it uses the Node interface
type Tree struct {
	root, focus Node
}

//Insert adds a window (and thus a leaf) to the tree.
//it adds it as a child of the currently focused node
//and then sets the focus to the newly added leaf
func (t *Tree) Insert(win *Window) {
	if t.focus == nil && t.root == nil {
		t.root = NewBranch()
		t.focus = t.root
	}
	l := NewLeaf(win)
	l.AttachToHook("leaf::remove", t.Remove)
	t.focus.AddChild(l)
	t.focus = l
}

//Remove removes the node passed to it.
//It is setup to work as a Signalable
//hook function. By convention the first
//argument should fulfil the Node interface
//any after that are ignored.
func (t *Tree) Remove(args ...interface{}) {
	n := args[0].(Node)
    nParent := n.Parent()
	if t.focus == n {
		t.focus = n.Parent()
	}
	nParent.RemoveChild(n)
    if len(nParent.Children()) == 0 && nParent.Parent() != nil{
        t.focus = nParent.Parent()
        nParent.Parent().RemoveChild(nParent)
    }
}

//Layout walks the tree calling the Layout
//method of each node
func (t *Tree) Layout(rect *xrect.XRect) {
	t.root.Layout(rect)
}

//FocusNext sets t.focus to the leaf after t.focus
func (t *Tree) FocusNext() {
	var l Node
	useNext := false
	t.root.VisitLeaves(func(n Node) bool {
		if useNext == true {
			l = n
			return false
		}
		if n == t.focus {
			useNext = true
		}
		return true
	})
	if l != nil {
		t.focus = l
	}
}

//FocusPrev sets t.focus to the leaf before t.focus
func (t *Tree) FocusPrev() {
	var l Node
	useNext := false
	t.root.VisitLeavesReverse(func(n Node) bool {
		if useNext == true {
			l = n
			return false
		}
		if n == t.focus {
			useNext = true
		}
		return true
	})
	if l != nil {
		t.focus = l
	}
}

//Print walks the tree printing each node
func (t *Tree) Print() {
	if t.root == nil {
		return
	}
	fmt.Println("Printing tree:")
	t.root.Print("")
	fmt.Println()
}
