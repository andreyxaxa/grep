package grep

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/andreyxaxa/pkg/grep/helpers"
)

// Params is core struct for grep util
type Params struct {
	reader io.ReadCloser
	target string

	// flags:
	A int
	B int
	C int

	n bool
	c bool
	i bool
	v bool
	F bool
}

// NewParams returns new 'Params' for grep-util
func NewParams() *Params {
	return &Params{}
}

// Start run grep
func (p *Params) Start() error {
	if err := p.parse(); err != nil {
		return err
	}

	defer p.reader.Close()

	lines := make([]string, 0)
	helpers.ReadLines(p.reader, &lines)

	p.Grep(lines)

	return nil
}

func (p *Params) parse() error {
	flag.BoolVar(&p.n, "n", false, "print line number for matched lines")
	flag.BoolVar(&p.c, "c", false, "print only a count of matching lines")
	flag.BoolVar(&p.i, "i", false, "print matched lines (ignore case)")
	flag.BoolVar(&p.F, "F", false, "fixed string, not a reg. exp")
	flag.BoolVar(&p.v, "v", false, "invert match: non-matching lines")
	flag.IntVar(&p.A, "A", 0, "N lines after match")
	flag.IntVar(&p.B, "B", 0, "N lines before match")
	flag.IntVar(&p.C, "C", 0, "N lines before and after match")

	flag.Parse()

	args := flag.Args()

	p.target = args[0]

	if p.C > 0 {
		p.A = p.C
		p.B = p.C
	}

	// смотрим, как поступают данные
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		p.reader = os.Stdin
		return nil
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open: %s - %v\n", os.Args[1], err)
		return err
	}
	p.reader = file

	return nil
}

// Grep is core function with main logic of grep util (Match() -> MakeRanges() -> MergeRanges() -> Print())
func (p *Params) Grep(lines []string) {
	matches := make([]int, 0)

	for i, line := range lines {
		if p.v {
			if !p.Match(line) {
				matches = append(matches, i)
			}
		} else {
			if p.Match(line) {
				matches = append(matches, i)
			}
		}
	}

	if p.c {
		fmt.Println(len(matches))
		return
	}

	if len(matches) == 0 {
		return
	}

	ranges := p.MakeRanges(matches, len(lines))

	merged := p.MergeRanges(ranges)

	for _, r := range merged {
		for i := r[0]; i <= r[1]; i++ {
			if p.n {
				fmt.Printf("%d:%s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
	}
}

// MakeRanges takes 'indexes'(lines with matches), 'ogLen'(len() of original text) and returns slice of ranges of matched lines
func (p *Params) MakeRanges(indexes []int, ogLen int) [][2]int {
	ranges := make([][2]int, 0)

	// например, если совпадение в строке 9, флаги: -A2, -B3, тогда
	// start = 9-3 = 6; end = 9+2 = 11;
	// Диапазон строк, которые нам нужно вывести - [6 11]
	for _, v := range indexes {
		start := v - p.B // начало
		if start < 0 {
			start = 0
		}
		end := v + p.A // конец
		if end >= ogLen {
			end = ogLen - 1
		}
		ranges = append(ranges, [2]int{start, end})
		/*
			[
				[6 11]
				[5 6]
				....
			]
		*/
	}

	return ranges
}

// MergeRanges takes 'ranges' from 'MakeRanges' and merges it if they intersect
func (p *Params) MergeRanges(ranges [][2]int) [][2]int {
	merged := make([][2]int, 0)

	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		r := ranges[i]
		// если начало следущего до конца текущего
		// если конец следущего после конца текущего
		// то сливаем
		if r[0] <= current[1]+1 {
			if r[1] > current[1] {
				current[1] = r[1]
			}
		} else {
			merged = append(merged, current)
			current = r
		}
	}
	merged = append(merged, current)

	return merged
}

// Match takes line and looking for matches (strings.Contains / regexp.MatchString)
func (p *Params) Match(line string) bool {
	t := p.target
	l := line
	var ok bool

	if p.i {
		t = strings.ToLower(t)
		l = strings.ToLower(l)
	}

	if p.F {
		ok = strings.Contains(l, t)
	} else {
		matched, err := regexp.MatchString(t, l)
		if err != nil {
			return false
		}
		ok = matched
	}

	return ok
}
