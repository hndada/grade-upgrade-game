package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func gacha(cdf []float64) int {
	r := rand.Float64()
	for i, v := range cdf {
		if r < v {
			return i
		}
	}
	return len(cdf) - 1
}

type Grade int // Level
//	func (g Grade) Output() int {
//		// +0, +1, +2, +3, +4, +5
//		// 0.54, 0.25, 0.125, 0.05, 0.02, 0.01, 0.005
//		cdf := []float64{.54, .79, .915, .965, .985, .995, 1.0}
//		return gacha(cdf)
//	}
func (g Grade) Output() float64 {
	cdf := []float64{.5, .9, 1.0}
	base := []float64{1, 18, 54}[gacha(cdf)]
	factor := 0.5 + rand.Float64()
	return base * factor
}
func (g Grade) Pass() bool {
	var minimum int
	if g <= 6 {
		minimum = 50
	} else if g <= 9 {
		minimum = 60
	} else if g <= 12 {
		minimum = 70
	} else if g <= 16 {
		minimum = 80
	} else if g <= 18 {
		minimum = 85
	} else if g <= 22 {
		minimum = 90
	} else {
		minimum = 95
	}
	return rand.Intn(100) >= minimum
}

type Book int // Level
func (b Book) Study() int {
	none := 0.33 + 0.01*float64(b)
	normal := 0.60 - 0.01*float64(b)
	book := 0.03 // - 0.005*float64(b)
	iq := 0.02   // + 0.005*float64(b)
	cdf := []float64{
		none,
		none + normal,
		none + normal + book,
		none + normal + book + iq,
		1.0,
	}
	return gacha(cdf)
}

// const (
// 	None     = iota
// 	Student  // 학생. 책 수준이 높으면 높은 등급의 학생이 나옴
// 	Semester // 학기. 4개 모이면 진학
// 	IQ       // 지능. 영구 능력치 상승
// )

// type Player struct {
// 	Grades   []Grade
// 	Semester int
// 	Book
// 	IQ int
// }

type Player struct {
	Grade
	Exp float64
	Book
	IQ int
}

// func NewPlayer() Player {
// 	p := Player{
// 		Grades:   make([]Grade, len(GradeNames)+5),
// 		Semester: 0,
// 		Book:     0,
// 		IQ:       0,
// 	}
// 	p.Grades[0] = 1
// 	return p
// }
// func (p Player) Grade() int {
// 	for i := len(p.Grades) - 1; i >= 0; i-- {
// 		if p.Grades[i] != 0 {
// 			return i
// 		}
// 	}
// 	return 0
// }

const (
	None = iota
	Normal
	UpgradeBook
	IQ
)

func (p *Player) Study() {
	switch p.Book.Study() {
	case None:
		// fmt.Print("허탕 쳤습니다. ")
	case Normal:
		// fmt.Print("공부 했습니다. ")
		p.Exp += p.Grade.Output()
	case UpgradeBook:
		if p.Book < Book(len(BookNames))-1 {
			// fmt.Print("책을 업그레이드 했습니다. ")
			p.Book++
		}
	case IQ:
		p.IQ++
		// fmt.Print("지능이 올랐습니다. ")
	}
}
func (p Player) Status() string {
	if p.IsClear() {
		return "축하합니다, 교수가 되셨습니다."
	}
	return fmt.Sprintf(
		"%s (%.2f%%). %s 공부중⋯. IQ는 %d \n",
		GradeNames[p.Grade], p.Exp, BookNames[p.Book], 100+p.IQ*5)
}

func (p Player) IsClear() bool { return p.Grade >= Grade(len(GradeNames)-1) }

//	func (p *Player) Test() {
//		for i, v := range p.Grades {
//			if v >= 2 {
//				if Grade(i).Pass() {
//					p.Grades[i] -= 2
//					p.Grades[i+1]++
//				} else {
//					p.Grades[i]--
//				}
//			}
//		}
//	}

func (p Player) Rest() {
	ipEffect := 10 * math.Log(float64(p.IQ+1))
	t := time.Second - time.Duration(ipEffect)*time.Millisecond - 750*time.Millisecond
	time.Sleep(t)
}
func main() {
	// p := NewPlayer()
	p := Player{}
	count := 0
	for !p.IsClear() {
		fmt.Printf("%d일째, %s", count, p.Status())
		p.Study()
		if p.Exp >= 100 {
			if p.Grade.Pass() {
				fmt.Println("합격!")
				p.Grade++
			} else {
				fmt.Println("불합ㅠㅠ")
			}
			p.Exp = 0
		}
		p.Rest()
		count++
	}
}