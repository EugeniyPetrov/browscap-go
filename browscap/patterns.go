package browscap

type Patterns []string

func patternLen(pattern string) int {
	l := 0
	for _, r := range pattern {
		if r == '*' || r == '?' {
			continue
		}
		l++
	}
	return l
}

func (p Patterns) Len() int {
	return len(p)
}

func (p Patterns) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Patterns) Less(i, j int) bool {
	pLen1 := patternLen(p[i])
	pLen2 := patternLen(p[j])

	if pLen1 != pLen2 {
		return pLen1 > pLen2
	}

	len1 := len(p[i])
	len2 := len(p[j])

	if len1 != len2 {
		return len1 > len2
	}

	return p[i] > p[j]
}
