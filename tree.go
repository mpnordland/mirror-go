package main

import (
    "github.com/BurntSushi/xgbutil/xrect"
)


type Node interface{
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
}

type NodeSlice []Node

func sliceRemove(ns NodeSlice, i int) NodeSlice{
    return append(ns[:i], ns[i+1:]...)
}

type Leaf struct{
    *Dispatcher
    parent Node
    win *Window
}

func NewLeaf(win *Window) *Leaf{
    l := &Leaf{NewDispatcher(), nil, win}
    win.AttachToHook("window::unmapped", l.removeSelf)
    return l
}

func (l *Leaf) AddChild(n Node){
    lParent := l.parent
    l.parent.RemoveChild(l)
    b := NewBranch()
    if lParent.SplitType() == "horizontal"{
        b.SplitTypeSet("vertical")
    }else{
        b.SplitTypeSet("horizontal")
    }
    b.AddChild(l)
    b.AddChild(n)
    lParent.AddChild(b)
    }
func (l *Leaf) RemoveChild(n Node){}
func (l *Leaf) SplitType() string {
    return ""
}
func (l *Leaf) SplitTypeSet(sType string){}

func (l *Leaf) Layout(rect *xrect.XRect){
    l.win.GeometrySet(rect)
}

func (l *Leaf) Parent() Node{
    return l.parent
}

func (l *Leaf) ParentSet(n Node){
    l.parent = n
}

func (l *Leaf) VisitLeaves(f func(n Node)bool) bool{
    return f(l)
}

func (l *Leaf) VisitLeavesReverse(f func(n Node)bool) bool{
    return f(l)
}

func (l *Leaf) removeSelf(args... interface{}){
    l.parent.RemoveChild(l)
}

func (l *Leaf) Children() NodeSlice{
    return make(NodeSlice, 0)
}

type Branch struct{
    *Dispatcher
    parent Node
    children NodeSlice
    splitType string
}

func NewBranch() *Branch{
    return &Branch{NewDispatcher(), nil, make(NodeSlice, 0), "horizontal"}
}

func (b *Branch) AddChild(n Node){
    n.ParentSet(b)
    b.children = append(b.children, n)
}

func (b *Branch) RemoveChild(n Node){
    for i, nn := range b.children{
        if nn == n{
            b.children = sliceRemove(b.children, i)
            b.children = append(b.children, n.Children()...)
            break
        }
    }
}

func (b *Branch) SplitType() string {
    return b.splitType
}

func (b *Branch) SplitTypeSet(sType string){
    b.splitType = sType
}

func (b *Branch) Layout(rect *xrect.XRect){
    if len(b.children) == 0{
        return
    }
    if b.splitType == "horizontal"{
        width := rect.Width()/len(b.children)
        x := rect.X()
        for _, n := range b.children{
            n.Layout(xrect.New(x, rect.Y(), width, rect.Height()))
            x += width
        }
    }else{
        height := rect.Height()/len(b.children)
        y := rect.Y()
        for _, n := range b.children{
            n.Layout(xrect.New(rect.X(), y, rect.Width(), height))
            y += height
        }
    }
}

func (b *Branch) Parent() Node{
    return b.parent
}

func (b *Branch) ParentSet(n Node){
    b.parent = n
}

func (b *Branch) VisitLeaves(f func(n Node) bool) bool{
    for _, c:= range b.children{
        if !c.VisitLeaves(f){
            return false
        }
    }
    return true
}

func (b *Branch) VisitLeavesReverse(f func(n Node) bool) bool{
    for i:=len(b.children)-1; i>=0; i--{
        c := b.children[i]
        if !c.VisitLeavesReverse(f){
            return false
        }
    }
    return true
}

func (b *Branch) Children() NodeSlice{
    return b.children
}

type Tree struct{
    root, focus Node
}

func (t *Tree) Insert(win *Window){
    if t.focus == nil && t.root == nil{
        t.root = NewBranch()
        t.focus = t.root
    }
    l := NewLeaf(win)
    t.focus.AddChild(l)
    t.focus = l
}

func (t *Tree) Remove(n Node){
    n.Parent().RemoveChild(n)
}

func (t *Tree) Layout(rect *xrect.XRect){
    t.root.Layout(rect)
}

func (t *Tree) FocusNext(){
    var l Node
    useNext := false
    t.root.VisitLeaves(func(n Node) bool{
        if useNext == true{
            l = n
            return false
        }
        if n == t.focus{
            useNext = true
        }
        return true
    })
    if l != nil{
        t.focus = l
    }
}

func (t *Tree) FocusPrev(){
    var l Node
    useNext := false
    t.root.VisitLeavesReverse(func(n Node) bool{
        if useNext == true{
            l = n
            return false
        }
        if n == t.focus{
            useNext = true
        }
        return true
    })
    if l != nil{
        t.focus = l
    }
}
