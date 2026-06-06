package domain

type Category string

const (
	CategoryPolicy        Category = "policy"
	CategoryAnswer        Category = "answer"
	CategoryClaim         Category = "claim"
	CategoryCertification Category = "certification"
	CategoryArchitecture  Category = "architecture"
)

func (c Category) Valid() bool {
	switch c {
	case CategoryPolicy, CategoryAnswer, CategoryClaim, CategoryCertification, CategoryArchitecture:
		return true
	}
	return false
}

func ValidCategories() []Category {
	return []Category{CategoryPolicy, CategoryAnswer, CategoryClaim, CategoryCertification, CategoryArchitecture}
}

func ParseCategory(s string) (Category, bool) {
	for _, c := range ValidCategories() {
		if string(c) == s {
			return c, true
		}
	}
	return "", false
}
