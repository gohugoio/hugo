// Code generated by "stringer -type SubCategory"; DO NOT EDIT.

package tplimpl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SubCategoryMain-0]
	_ = x[SubCategoryEmbedded-1]
	_ = x[SubCategoryInline-2]
}

const _SubCategory_name = "SubCategoryMainSubCategoryEmbeddedSubCategoryInline"

var _SubCategory_index = [...]uint8{0, 15, 34, 51}

func (i SubCategory) String() string {
	if i < 0 || i >= SubCategory(len(_SubCategory_index)-1) {
		return "SubCategory(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SubCategory_name[_SubCategory_index[i]:_SubCategory_index[i+1]]
}
