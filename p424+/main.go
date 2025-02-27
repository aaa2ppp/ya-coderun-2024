package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func run(in io.Reader, out io.Writer) error {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	num, err := io.ReadAll(br)
	if err != nil {
		return err
	}
	num = bytes.TrimSpace(num)
	n := len(num)
	if n&1 != 0 {
		return fmt.Errorf("%d is not even", n)
	}

	if _, flag := addOne(0, num, len(num)-1); flag {
		addOne(0, num, len(num)-1)
	}

	part1 := num[:n/2]
	part2 := num[n/2:]
	sum1 := calcSum(part1)
	sum2 := calcSum(part2)

	sum2, flag := rotate1(sum1, sum2, part2)
	if flag {
		sum1, _ = addOne(sum1, part1, len(part1)-1)
	}
	rotate2(sum1, sum2, part2)

	bw.Write(num)
	bw.WriteByte('\n')

	return nil
}

// суммирует все цифры числа
func calcSum(num []byte) int {
	sum := 0
	for _, v := range num {
		sum += int(v - '0')
	}
	return sum
}

// добавляет 1 в позицию i (отсчет от 0 слева направо), возвращает новую сумму и флаг переполнения
func addOne(sum int, num []byte, i int) (int, bool) {
	for i >= 0 && num[i] == '9' {
		sum -= 9
		num[i] = '0'
		i--
	}

	if i == -1 {
		return sum, true
	}

	sum++
	num[i]++
	return sum, false
}

// крутит вторую часть, пока она больше первой, возвращает новую сумму второй части и флаг переполнения
func rotate1(sum1, sum2 int, part2 []byte) (int, bool) {
	i := len(part2) - 1

	for i >= 0 && sum2-sum1 > 0 {
		sum2 -= int(part2[i] - '0')
		part2[i] = '0'

		var flag bool
		sum2, flag = addOne(sum2, part2, i-1)
		if flag {
			return sum2, true
		}

		i--
	}

	if i == -1 {
		return sum2, true
	}

	return sum2, false
}

// крутит вторую часть, пока она меньше первой
func rotate2(sum1, sum2 int, part2 []byte) {
	diff := sum1 - sum2
	if diff < 0 {
		log.Panicf("sum1 must be > mum2; sum1:%d sum2:%d part2:%v", sum1, sum2, part2)
	}
	i := len(part2) - 1
	for diff > int('9'-part2[i]) {
		diff -= int('9' - part2[i])
		part2[i] = '9'
		i--
	}
	part2[i] += byte(diff)
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}