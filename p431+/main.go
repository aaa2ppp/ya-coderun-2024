package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type point struct {
	x, y int
}

type segment struct {
	p0, p1 point
}

// возвращает true, если точка пренажлежит прямой отрезка
func (s segment) contains(p point) bool {
	return (p.x-s.p0.x)*(s.p1.y-s.p0.y) == (p.y-s.p0.y)*(s.p1.x-s.p0.x)
}

// возвращает "растояние" между отрезком и точкой (умноженное на длину отрезка?)
// нам этого должно хватитить для сравнения удаленности точек от торезка
func (s segment) distance(p point) int {
	v := abs((s.p1.y-s.p0.y)*p.x - (s.p1.x-s.p0.x)*p.y + s.p1.x*s.p0.y - s.p1.y*s.p0.x)
	return v
}

// возвращает индекс наиболие дяльней точки к прямой отрезка
func (s segment) seachFarthest(pp []point) int {
	var i2 int
	maximum := s.distance(pp[i2])
	for i := range pp {
		if d := s.distance(pp[i]); d > maximum {
			maximum = d
			i2 = i
		}
	}
	return i2
}

// возвращает индекс ближайшей точки к прямой отрезка
func (s segment) seachNearest(pp []point) int {
	var i3 int
	minimum := s.distance(pp[i3])
	for i := range pp {
		if d := s.distance(pp[i]); d < minimum {
			minimum = d
			i3 = i
		}
	}
	return i3
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// возвращает true, если отрезки прересекаются. если конец одного отрезка пренадлежит другому,
// это не считается пересечением
func intersect(s1, s2 segment) bool {
	v1 := s1.p1.x - s1.p0.x
	w1 := s1.p1.y - s1.p0.y
	v2 := s2.p1.x - s2.p0.x
	w2 := s2.p1.y - s2.p0.y
	v3 := s2.p0.x - s1.p0.x
	w3 := s2.p0.y - s1.p0.y

	z1 := w1*v3 - v1*w3
	z2 := w2*v3 - v2*w3
	z3 := v1*w2 - w1*v2

	sign := func(a int) int {
		if a < 0 {
			return -1
		}
		if a > 0 {
			return 1
		}
		return 0
	}

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	check := func(u1, u2 int) bool {
		if sign(u1) != sign(u2) {
			return false
		}
		return abs(u1) < abs(u2)
	}

	return check(z1, z3) && check(z2, z3)
}

func solve(pp []point) bool {
	if len(pp) <= 6 {
		return slowSolve(pp)
	}

	n := len(pp)
	_ = n

	removePoint := func(i int) {
		n--
		pp[i] = pp[n]
		pp = pp[:n]
	}

	// стоим "шесть отрезков" по четырем точкам учитывая нюансы
	ss := make([]segment, 0, 6)

	// будем считать внутренние точки отрезков и количество отрезков
	// содедержащих внутренние точки
	holder := make([]int, 0, 6)
	holderCount := 0

	// ищем наиболее удаленные друг от друга точки по х
	var p0, p1 point
	{
		var i0, i1 int
		for i := range pp {
			if pp[i].x < pp[i0].x || pp[i].x == pp[i0].x && pp[i].y < pp[i0].y {
				i0 = i
			}
			if pp[i].x > pp[i1].x || pp[i].x == pp[i1].x && pp[i].y > pp[i1].y {
				i1 = i
			}
		}

		p0 = pp[i0]
		p1 = pp[i1]
		if debugEnable {
			log.Println("found p0, p1:", p0, p1)
		}

		if i0 < i1 { // xxx
			i0, i1 = i1, i0
		}
		removePoint(i0)
		removePoint(i1)
		if debugEnable {
			log.Println("pp:", pp)
		}
	}

	ss = append(ss, segment{p0, p1})
	holder = append(holder, 0)

	// ищем из оставшихся наиболее удаленную точку от отрезка
	var p2 point
	{
		i2 := ss[0].seachFarthest(pp)
		p2 = pp[i2]
		if debugEnable {
			log.Println("found p2:", p2)
		}

		removePoint(i2)
		if debugEnable {
			log.Println("pp:", pp)
		}
	}

	if ss[0].contains(p2) {
		if debugEnable {
			log.Println("on one line")
		}
		return false
	}

	ss = append(ss, segment{p2, p0}, segment{p2, p1})
	holder = append(holder, 0, 0)

	// <KZNM!!!

	var p3 point
	{
		var i3 int
		if n == 1 {
			// только 4 точки - просто пересечение ищем. можно вынести наверх
			p3 = pp[0]
			removePoint(0)
			ss = append(ss, segment{p3, p0}, segment{p3, p1}, segment{p3, p2})
			holder = append(holder, 0, 0, 0)
		} else if ss[0].contains(pp[0]) && ss[0].contains(pp[1]) {
			holder[0] = 2
			holderCount++
			removePoint(1)
			removePoint(0)
		} else if ss[1].contains(pp[0]) && ss[1].contains(pp[1]) {
			holder[1] = 2
			holderCount++
			removePoint(1)
			removePoint(0)
		} else if ss[2].contains(pp[0]) && ss[2].contains(pp[1]) {
			holder[2] = 2
			holderCount++
			removePoint(1)
			removePoint(0)
		} else {
			s := segment{pp[0], pp[1]}
			if s.contains(p0) {
				i3 = ss[2].seachNearest(pp)
			} else if s.contains(p1) {
				i3 = ss[1].seachNearest(pp)
			} else if s.contains(p2) {
				i3 = ss[0].seachNearest(pp)
			} else {
				if debugEnable {
					log.Println("xxx")
				}
				return true
			}

			p3 = pp[i3]
			if debugEnable {
				log.Println("found p3:", p3)
			}

			removePoint(i3)
			if debugEnable {
				log.Println("pp:", pp)
			}

			// <KZNM!!!

			if ss[0].contains(p3) {
				ss = append(ss, segment{p3, p2})
				holder = append(holder, 1)
				holderCount++
			} else if ss[1].contains(p3) {
				ss = append(ss, segment{p3, p1})
				holder = append(holder, 1)
				holderCount++
			} else if ss[2].contains(p3) {
				ss = append(ss, segment{p3, p0})
				holder = append(holder, 1)
				holderCount++
			} else {
				ss = append(ss, segment{p3, p0}, segment{p3, p1}, segment{p3, p2})
				holder = append(holder, 0, 0, 0)
			}
		}
	}

	if debugEnable {
		log.Println("ss:", ss)
	}

	// ищем пересечения отрезков, если таковые имеются, то p0,p1,p2,p3 образуют выпуклый четырехугольник
	for i := range ss {
		for j := i + 1; j < len(ss); j++ {
			if intersect(ss[i], ss[j]) {
				if debugEnable {
					log.Println("intersect", ss[i], ss[j])
				}
				return true
			}
		}
	}

	// остальные точки являются внутренними. выпуклый четырехугольник не может быть построен,
	// только если все они принадлежат одному и тому же отрезку

	for _, p := range pp {
		if debugEnable {
			log.Println("+", p)
		}

		freePoint := true
		for j := range ss {
			if ss[j].contains(p) {
				freePoint = false
				holder[j]++
				if holder[j] == 1 {
					holderCount++
					if holderCount > 1 {
						if debugEnable {
							log.Println("holder count > 1")
						}
						return true
					}
				}

				if !(min(ss[j].p0.x, ss[j].p1.x) <= p.x && p.x <= max(ss[j].p0.x, ss[j].p1.x) &&
					min(ss[j].p0.y, ss[j].p1.y) <= p.y && p.y <= max(ss[j].p0.y, ss[j].p1.y)) {
					return true
				}
			}
		}

		if freePoint {
			if debugEnable {
				log.Println("free point (outside any segment)")
			}
			return true
		}
	}

	return false
}

func slowSolve(pp []point) bool {
	for i := range pp {
		p0 := pp[i]
		for j := i + 1; j < len(pp); j++ {
			p1 := pp[j]
			for k := j + 1; k < len(pp); k++ {
				p2 := pp[k]
				for l := k + 1; l < len(pp); l++ {
					p3 := pp[l]
					if intersect(segment{p0, p1}, segment{p2, p3}) ||
						intersect(segment{p0, p3}, segment{p1, p2}) ||
						intersect(segment{p0, p2}, segment{p1, p3}) {
						return true
					}
				}
			}
		}
	}
	return false
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	pp := make([]point, n)
	for i := range pp {
		x, y, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		pp[i] = point{x: x, y: y}
	}

	res := solve(pp)

	if res {
		bw.WriteString("Yes\n")
	} else {
		bw.WriteString("No\n")
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
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

// ----------------------------------------------------------------------------

func gcd(a, b int) int {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
