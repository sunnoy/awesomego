package le09_02

// nums []int{2, 7, 11, 15},
//target 17
func TwoSum(nums []int, target int) []int {

	m := make(map[int]int)

	for i := 0; i < len(nums); i++ {

		a := target - nums[i]

		if _, ok := m[a]; ok {
			return []int{m[a], i}
		}
		m[nums[i]] = i

	}

	return nil

}
