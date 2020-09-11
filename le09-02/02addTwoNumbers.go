package le09_02

//2 个逆序的链表，要求从低位开始相加，得出结果也逆序输出，返回值是逆序结果链表的头结点。
type ListNode struct {
	Val  int
	Next *ListNode
}

//l1 := &le09_02.ListNode{
//	Val: 2,
//	Next: &le09_02.ListNode{
//		Val: 4,
//		Next: &le09_02.ListNode{
//			Val: 3,
//		},
//	},
//}
//
//l2 := &le09_02.ListNode{
//	Val: 5,
//	Next: &le09_02.ListNode{
//		Val: 6,
//		Next: &le09_02.ListNode{
//			Val: 4,
//		},
//	},
//}

//为了处理方法统一，可以先建立一个虚拟头结点，这个虚拟头结点的 Next 指向真正的 head，这样 head 不需要单独处理，直接 while 循环即可。
//另外判断循环终止的条件不用是 p.Next ！= nil，这样最后一位还需要额外计算，循环终止条件应该是 p != nil。
func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	// 创建虚拟头节点
	head := &ListNode{Val: 0}

	n1, n2, carry, current := 0, 0, 0, head

	// l1 或 l2 都不为空，carry 不为0的时候
	for l1 != nil || l2 != nil || carry != 0 {

		// 如果l1 为空
		if l1 == nil {
			n1 = 0
		} else {
			n1 = l1.Val
			l1 = l1.Next

		}

		if l2 == nil {
			n2 = 0

		} else {
			n2 = l2.Val
			l2 = l2.Next
		}

		current.Next = &ListNode{
			Val: (n1 + n2 + carry) % 10,
		}
		current = current.Next

		carry = (n1 + n2 + carry) / 10

	}
	return head.Next

}
