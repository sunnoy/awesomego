package main

import (
	"fmt"
	le09_02 "leetcode/le09-02"
)

func main() {
	//fmt.Printf("%v", le09_02.TwoSum([]int{2, 7, 11, 15}, 17))

	l1 := &le09_02.ListNode{
		Val: 2,
		Next: &le09_02.ListNode{
			Val: 4,
			Next: &le09_02.ListNode{
				Val: 3,
			},
		},
	}

	l2 := &le09_02.ListNode{
		Val: 5,
		Next: &le09_02.ListNode{
			Val: 6,
			Next: &le09_02.ListNode{
				Val: 4,
			},
		},
	}

	listnode := le09_02.AddTwoNumbers(l1, l2)

	fmt.Printf("%v\n %v\n %v", listnode.Val, listnode.Next.Val, listnode.Next.Next.Val)

}
