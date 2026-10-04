package homework01

import (
	"fmt"
	"sort"
	"strconv"
)

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	counts := make(map[int]int)

	for _, n := range nums {
		counts[n]++
	}

	for k, v := range counts {
		if v == 1 {
			return k
		}
	}

	return 0
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	}

	if x < 10 {
		return true
	}

	if x%10 == 0 && x > 10 {
		return false
	}

	s := strconv.Itoa(x)

	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		if s[l] != s[r] {
			return false
		}
	}

	return true
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	stack := []byte{}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if m, ok := pairs[c]; ok {
			//这里说明是右括号，出栈
			// 【防御】如果是个右括号，但栈是空的，说明它没有对应的左括号，空栈再执行出线的操作会导致panic: runtime error: index out of range [-1]
			if len(stack) == 0 {
				fmt.Printf("Invalid: 空栈情况下先出现了右括号%c\n", c)
				return false
			}
			//取栈顶
			top := stack[len(stack)-1]
			if top == m {
				stack = stack[:len(stack)-1]
			} else {
				fmt.Printf("Invalid: 出现的右括号%c与栈顶的左括号%c不配对\n", c, top)
				return false
			}
		} else {
			//持续出现左括号，入栈
			stack = append(stack, c)
			fmt.Println(string(stack))
		}
	}

	if len(stack) == 0 {
		return true
	}

	fmt.Printf("Invalid: 栈中残余左括号%s\n", string(stack))

	return false
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	for i := 0; i < len(strs[0]); i++ {
		for j := 1; j < len(strs); j++ {

			if len(strs[j]) == i {
				fmt.Printf("%s\n", strs[j])
				return strs[0][:i]
			}

			if strs[0][i] != strs[j][i] {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	//大端模式存储，默认小学算数进位答题
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		} else {
			digits[i] = 0
		}
	}

	result := make([]int, len(digits)+1)
	result[0] = 1
	return result
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	l := len(nums)
	if l == 0 {
		return 0
	}

	i := 0

	for j := i + 1; j < l; j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}

	fmt.Printf("%v\n", i)

	return i + 1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	sort.Slice(
		intervals,
		func(i, j int) bool {
			// 按照区间的左边界（第 0 个元素）从小到大排序
			return intervals[i][0] < intervals[j][0]
		})
	fmt.Printf("%v\n", intervals)
	merged := [][]int{}

	merged = append(merged, intervals[0])
	prev := 0

	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] > merged[prev][1] {
			merged = append(merged, intervals[i])
			prev++
		} else if intervals[i][1] > merged[prev][1] {
			merged[prev][1] = intervals[i][1]
		}
	}
	fmt.Printf("Merged: %v\n", merged)

	return merged
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	// TODO: implement
	return nil
}
