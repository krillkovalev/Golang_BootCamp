// Создается стек, начиная с текущего узла (*n).
// Пока стек не пустой:
// Извлекаем последний элемент стека (текущий узел).
// Если имя текущего узла соответствует искомому, возвращаем его.
// Для каждого дочернего узла добавляем его в конец стека.
// Если искомый узел не найден, возвращается nil.

// func main() {
// 	a := &TreeNode{HasToy: false}
// 	b := &TreeNode{HasToy: true}
// 	c := &TreeNode{HasToy: false}
// 	e := &TreeNode{HasToy: true}
// 	g := &TreeNode{HasToy: true}
// 	a.Left = b
// 	a.Right = c
// 	a.Left.Right = e 
// 	a.Right.Right = g
	

// 	areToysBalanced(a)

// }

// func main() {
// 	a := &TreeNode{HasToy: false}
// 	b := &TreeNode{HasToy: false}
// 	c := &TreeNode{HasToy: true}
// 	a.Left = b
// 	a.Right = c
// 	d := &TreeNode{HasToy: false}
// 	e := &TreeNode{HasToy: true}
// 	a.Left.Left = d
// 	b.Left.Right = e

// 	areToysBalanced(a)

// }

// func main() {
// 	a := &TreeNode{HasToy: true}
// 	b := &TreeNode{HasToy: true}
// 	c := &TreeNode{HasToy: false}
// 	a.Left = b
// 	a.Right = c
// 	d := &TreeNode{HasToy: true}
// 	e := &TreeNode{HasToy: false}
// 	f := &TreeNode{HasToy: true}
// 	g := &TreeNode{HasToy: true}
// 	a.Left.Left = d
// 	a.Left.Right = e
// 	a.Right.Left = f
// 	a.Right.Right = g

// 	areToysBalanced(a)

// }




package main 

import (
	"fmt"
)

type TreeNode struct {
	HasToy bool
	Left *TreeNode
	Right *TreeNode
}


func main() {
	a := &TreeNode{HasToy: true}
	b := &TreeNode{HasToy: true}
	c := &TreeNode{HasToy: false}
	a.Left = b
	a.Right = c
	d := &TreeNode{HasToy: true}
	e := &TreeNode{HasToy: false}
	f := &TreeNode{HasToy: true}
	g := &TreeNode{HasToy: true}
	a.Left.Left = d
	a.Left.Right = e
	a.Right.Left = f
	a.Right.Right = g

	fmt.Println(areToysBalanced(a))

}



func pre_order_traversal(tree *TreeNode) int{
	count := 0
	if tree == nil {
		return 0
	}
	if tree.HasToy {
		count += 1
	}

	count += pre_order_traversal(tree.Left)
	count += pre_order_traversal(tree.Right)

	return count
}

func areToysBalanced(tree *TreeNode) bool{
	if tree.Left == nil && tree.Right == nil {
		return true
	}

	left_subtree := pre_order_traversal(tree.Left)
	right_subtree := pre_order_traversal(tree.Right)

	if left_subtree != right_subtree {
		return false
	}

	return true 

}

func unrollGarland() []bool{

}