package multiverse

func split(input string) []string {
	var s []string
	start := 0
	end := 0
	l := len(input)
	size := 1

	for i := 0; i < l; i++ {
		if input[i] == ' ' {
			size++
		}
	}

	if size == 1 {
		return []string{input}
	}

	s = make([]string, size)

	for i := 0; i < size; i++ {
		for end < l && input[end] != ' ' {
			end++
		}

		s[i] = input[start:end]

		end++
		start = end
	}

	return s
}
