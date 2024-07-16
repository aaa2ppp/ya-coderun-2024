package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type Item struct {
	red_idx       int
	next_blue_idx int
	next_offset   int
	index         int // The index of the item in the heap.
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].next_offset < pq[j].next_offset
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) { /* stub */ }
func (pq *PriorityQueue) Pop() any   { /* stub */ return nil }

func calcIntesection(blue, red []int) int {
	if debugEnable {
		log.Println("-- calcIntesection")
		log.Println("blue:", blue)
		log.Println("red: ", red)
	}
	intersection := 0

	for i, j := 1, 1; i < len(blue) && j < len(red); {
		a, b := blue[i-1], blue[i]
		c, d := red[j-1], red[j]
		if debugEnable {
			log.Println("a,b:", a, b)
			log.Println("c,d:", c, d)
		}

		if b < c {
			i += 2
			continue
		}
		if d < a {
			j += 2
			continue
		}

		intersection += min(b, d) - max(a, c)

		if b < d {
			i += 2
		} else {
			j += 2
		}
	}

	return intersection
}

// функция считает, что первая и последняя синие точки являются общими границами
func prepareItems(blue, red []int) []Item {
	if debugEnable {
		log.Println("-- prepareItems")
		log.Println("blue:", blue)
		log.Println("red: ", red)
	}
	items := make([]Item, len(red))

	i, j := 1, 0
	for i < len(blue) && j < len(red) {
		b := blue[i]
		c := red[j]

		if b <= c {
			i++
			continue
		}

		if debugEnable {
			log.Println("a,b,c:", blue[i-1], b, c)
		}

		items[j] = Item{
			red_idx:       j,
			next_blue_idx: i,
			next_offset:   b - c,
		}

		j++
	}

	for j < len(red) {
		items[j] = Item{
			red_idx:       j,
			next_blue_idx: len(blue),
		}
		j++
	}

	return items
}

func getSegmentAdding(c, d *Item) int {
	if c.next_blue_idx%2 == 1 && d.next_blue_idx%2 == 0 {
		// начало красного отрезка вне синего, а конец внутри
		return 1
	}
	if c.next_blue_idx%2 == 0 && d.next_blue_idx%2 == 1 {
		// начало красного отрезка внутри синего, а конец вне
		return -1
	}
	return 0
}

func caclAdding(items []Item) int {
	adding := 0

	for i := 1; i < len(items); i += 2 {
		c, d := &items[i-1], &items[i]
		if debugEnable {
			log.Println("getSegmentAdding:", getSegmentAdding(c, d))
		}
		adding += getSegmentAdding(c, d)
	}

	return adding
}

func moveTo(red []int, pos int) {
	offset := pos - red[0]
	for i := range red {
		red[i] += offset
	}
}

// функция считает, что первая и последняя синие точки являются общими границами
func solve(l, r int, blue_segs, red_segs [][2]int) int {
	n := len(blue_segs)
	m := len(red_segs)

	buf := make([]int, n*2+2) // резервируем 2 элемента под границы
	blue := buf[1 : len(buf)-1]
	sortSegments(blue_segs)
	if debugEnable {
		log.Println("blue:", blue_segs)
	}
	nn := getSegmentPoints(blue_segs, blue)
	if nn == 0 {
		return 0
	}
	blue = buf[:nn+2]
	blue[0] = l
	blue[len(blue)-1] = r
	if debugEnable {
		log.Println("blue:", blue)
	}

	red := make([]int, m*2)
	sortSegments(red_segs)
	if debugEnable {
		log.Println("red:", red_segs)
	}
	mm := getSegmentPoints(red_segs, red)
	if mm == 0 {
		return 0
	}
	red = red[:mm]
	if debugEnable {
		log.Println("red:", red)
	}

	moveTo(red, blue[0])
	if debugEnable {
		log.Println("red:", red)
	}

	intersection := calcIntesection(blue[1:len(blue)-1], red)
	if debugEnable {
		log.Println("intersection:", intersection)
	}

	items := prepareItems(blue, red)
	adding := caclAdding(items)
	if debugEnable {
		log.Println("items: ", items)
		log.Println("adding:", adding)
	}

	if debugEnable {
		log.Println("-- search minimum")
	}
	pq := make(PriorityQueue, len(items))
	for i := range pq {
		items[i].index = i
		pq[i] = &items[i]
	}
	heap.Init(&pq)

	offset := 0
	min_intersection := intersection

	move_next := func(it *Item) {
		it.next_blue_idx++
		if i := it.next_blue_idx; i < len(blue) {
			a, b := blue[i-1], blue[i]
			it.next_offset = offset + (b - a)
		} else {
			it.next_offset = 0
		}
		heap.Fix(&pq, it.index)
	}

	for {
		it := pq[0]
		if it.next_blue_idx == len(blue) {
			break
		}

		intersection += (it.next_offset - offset) * adding
		min_intersection = min(min_intersection, intersection)
		offset = it.next_offset

		var c, d *Item
		if it.red_idx%2 == 0 {
			c = it
			d = &items[it.red_idx+1]
		} else {
			c = &items[it.red_idx-1]
			d = it
		}

		if debugEnable {
			log.Println("getSegmentAdding:", getSegmentAdding(c, d))
		}
		adding -= getSegmentAdding(c, d)
		if c.next_offset == d.next_offset {
			move_next(c)
			move_next(d)
		} else {
			move_next(it)
		}
		if debugEnable {
			log.Println("getSegmentAdding:", getSegmentAdding(c, d))
		}
		adding += getSegmentAdding(c, d)

		if debugEnable {
			log.Println("items: ", items)
			log.Println("adding:", adding)
		}
	}

	return min_intersection
}

func sortSegments(segments [][2]int) {
	sort.Slice(segments, func(i, j int) bool {
		return segments[i][0] < segments[j][0] || segments[i][0] == segments[j][0] && segments[i][1] < segments[j][1]
	})

	// проверяем нет ли подлянки
	if checkIntersections(segments) {
		panic("there are intersecting segments")
	}
}

// Принимает слайс отсортированных по возрастанию первой точки отрезков.
// Возвращает истину если есть пересекающиеся одна общая точка не считается
// пересечением.
func checkIntersections(segments [][2]int) bool {
	for i := 1; i < len(segments); i++ {
		a, b := segments[i-1][0], segments[i-1][1]
		c, d := segments[i][0], segments[i][1]
		if a < d && c < b {
			return true
		}
	}
	return false
}

// Принимает слайс отсортированных по возрастанию первой точки отрезков.
// Записывает точки отрезков в слайс points пропуская дубли.
// Возвращает количество записанных точек.
func getSegmentPoints(segments [][2]int, points []int) int {
	nn := 0
	points = points[:0]
	for i := range segments {
		a, b := segments[i][0], segments[i][1]
		// if a == b {
		// 	continue
		// }
		if nn > 0 && a == points[nn-1] {
			points[nn-1] = b
			continue
		}
		nn += 2
		points = append(points, a, b)
	}

	return nn
}

func scanSegments(sc *bufio.Scanner, n int) ([][2]int, error) {
	segments := make([][2]int, n)
	for i := 0; i < n; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			return nil, err
		}
		segments[i] = [2]int{a, b}
	}
	return segments, nil
}

func run(in io.Reader, out io.Writer) (err error) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, l, r, err := scanFourInt(sc)
	if err != nil {
		return err
	}

	blue, err := scanSegments(sc, n)
	if err != nil {
		return err
	}

	red, err := scanSegments(sc, m)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()

	writeInt(bw, solve(l, r, blue, red), defaultWriteOpts)
	return nil
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanInt(sc *bufio.Scanner) (int, error) {
	sc.Scan()
	return strconv.Atoi(unsafeString(sc.Bytes()))
}

func scanTwoInt(sc *bufio.Scanner) (v1, v2 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	return v1, v2, err
}

func scanThreeInt(sc *bufio.Scanner) (v1, v2, v3 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	return v1, v2, v3, err
}

func scanFourInt(sc *bufio.Scanner) (v1, v2, v3, v4 int, err error) {
	v1, err = scanInt(sc)
	if err == nil {
		v2, err = scanInt(sc)
	}
	if err == nil {
		v3, err = scanInt(sc)
	}
	if err == nil {
		v4, err = scanInt(sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts(sc *bufio.Scanner, a []int) error {
	for i := range a {
		v, err := scanInt(sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type writeOpts struct {
	sep byte
	end byte
}

var defaultWriteOpts = writeOpts{
	sep: ' ',
	end: '\n',
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	_, err := bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))

	if err == nil && opts.end != 0 {
		bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
