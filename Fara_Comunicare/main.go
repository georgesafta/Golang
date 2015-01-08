// Foca_fara_Comunicare project main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var k = make([]int, 8)
var terminoiu = make(chan string)
var mutex sync.Mutex
var mutexfin sync.Mutex
var mutex1 sync.Mutex
var mutex2 sync.Mutex

type Robot struct {
	id       int
	obj_ch   int
	flag_ch  bool
	mp2      int // cel mai mic obiect
	nr_poste int // postul la care se afla
	tip      int // 13 sau 32
}

type Poste struct {
	nr  int
	obj []int
}

type Info struct {
	id          int // id robot
	post        int // post destinatie
	nr_incarcat int // id obiect incarcat
}

func ch_obj(p *Poste, r *Robot) { // incarcarea primului obiect din lista
	mutex.Lock()
	if len(p.obj) != 0 {
		r.flag_ch = true
		r.obj_ch = p.obj[0]
	}
	if len(p.obj) > 1 {
		p.obj = p.obj[1:]
	} else {
		p.obj = []int{}
	}
	if r.mp2 == 30000 {
		r.mp2 = r.obj_ch
	}
	if r.flag_ch == true {
		fmt.Printf("Incarcare robot #%d cu : %d de la postul %d\n", r.id, r.obj_ch, p.nr)
	}
	if len(p.obj) > 0 {
		continut := "Postul " + strconv.Itoa(p.nr) + " contine:"
		for i := 0; i < len(p.obj); i++ {
			//fmt.Printf("%d ", p.obj[i])
			continut += " " + strconv.Itoa(p.obj[i])
		}
		continut += "\n"
		fmt.Print(continut)
	} else {
		fmt.Printf("Postul %d este gol\n", p.nr)
	}
	mutex.Unlock()
}

func ch_obj_fin(p *Poste, r *Robot) { // incarcarea ultimlui obiect din lista
	mutexfin.Lock()
	f := 0
	if len(p.obj) != 0 {
		r.flag_ch = true
		f = len(p.obj) - 1
		r.obj_ch = p.obj[f]
	}
	if len(p.obj) > 1 {
		p.obj = p.obj[0:f]
	} else {
		p.obj = []int{}
	}
	if r.flag_ch == true {
		fmt.Printf("Incarcare robot #%d cu : %d de la postul %d\n", r.id, r.obj_ch, p.nr)
	}
	if len(p.obj) > 0 {
		continut := "Postul " + strconv.Itoa(p.nr) + " contine:"
		for i := 0; i < len(p.obj); i++ {
			//fmt.Printf("%d ", p.obj[i])
			continut += " " + strconv.Itoa(p.obj[i])
		}
		continut += "\n"
		fmt.Print(continut)
	} else {
		fmt.Printf("Postul %d este gol\n", p.nr)
	}
	mutexfin.Unlock()
}

func dch_obj(p *Poste, r *Robot) { // descarcare obiect
	//mutex.Lock()
	p.obj = append(p.obj, r.obj_ch)
	fmt.Printf("Robotul #%d descarca : %d la postul %d\n", r.id, r.obj_ch, p.nr)
	if len(p.obj) > 0 {
		continut := "Postul " + strconv.Itoa(p.nr) + " contine:"
		for i := 0; i < len(p.obj); i++ {
			//fmt.Printf("%d ", p.obj[i])
			continut += " " + strconv.Itoa(p.obj[i])
		}
		continut += "\n"
		fmt.Print(continut)
	} else {
		fmt.Printf("Postul %d este gol\n", p.nr)
	}
	//mutex.Unlock()
	r.flag_ch = false
	r.obj_ch = 30000
}

func vasPoste(p *Poste, r *Robot) {
	x := rand.Float32()
	switch r.tip {
	case 13:
		switch p.nr {
		case 1:
			switch {
			case x < .15:
				r.nr_poste = 3
			default:
				r.nr_poste = 1
			}
		case 3:
			switch {
			case x < .15:
				r.nr_poste = 1
			default:
				r.nr_poste = 3
			}
		}
	case 32:
		switch p.nr {
		case 2:
			switch {
			case x < .15:
				r.nr_poste = 3
			default:
				r.nr_poste = 2
			}
		case 3:
			switch {
			case x < .15:
				r.nr_poste = 2
			default:
				r.nr_poste = 3
			}
		}
	}
	fmt.Printf("Robotul #%d merge la postul : %d %f\n", r.id, r.nr_poste, x)
}

func logique_1(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, s []int) {
	if r.flag_ch == true {
		vasPoste(p3, r)
		logique(p1, p2, p3, r, s)
	} else {
		if len(p1.obj) != 0 {
			ch_obj(p1, r)
			vasPoste(p3, r)
			logique(p1, p2, p3, r, s)
		} else {
			fmt.Printf("Robotul #%d se opreste \n", r.id)
		}
	}
}

func logique_2(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, s []int) {
	mutex1.Lock()
	f := len(p2.obj)
	if r.flag_ch == true {
		switch {
		case f != 0:
			{
				if r.obj_ch > p2.obj[len(p2.obj)-1] {
					dch_obj(p2, r)
					if len(p3.obj) != 0 {
						if p2.obj[len(p2.obj)-1] > p3.obj[0] && r.flag_ch == false {
							ch_obj_fin(p2, r)
							if r.obj_ch < r.mp2 {
								r.mp2 = r.obj_ch
							}
						}
					}
				}
				if stop(s, p2) == false {
					vasPoste(p3, r)
				}
			}
		case f == 0:
			{
				dch_obj(p2, r)
				if len(p3.obj) != 0 {
					if p2.obj[len(p2.obj)-1] > p3.obj[0] {
						ch_obj_fin(p2, r)
						if r.obj_ch < r.mp2 {
							r.mp2 = r.obj_ch
						}
					}
				}
				if stop(s, p2) == false {
					vasPoste(p3, r)
				}
			}
		}
	} else {
		if f != 0 {
			if len(p3.obj) != 0 {
				if p2.obj[len(p2.obj)-1] > p3.obj[0] {
					ch_obj_fin(p2, r)
					if r.obj_ch < r.mp2 {
						r.mp2 = r.obj_ch
					}
				}
				vasPoste(p3, r)
			}
		}
	}
	mutex1.Unlock()
	if stop(s, p2) == false {
		logique(p1, p2, p3, r, s)
	}
}

func logique_3(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, s []int) {
	switch r.tip {
	case 13:
		if r.flag_ch == true {
			dch_obj(p3, r)
		}
		vasPoste(p1, r)
	case 32:
		mutex2.Lock()
		f3 := len(p3.obj)
		f2 := len(p2.obj)
		switch f3 {
		case 0:
			if r.flag_ch == true {
				if f2 == 0 {
					vasPoste(p2, r)
				} else {
					if r.obj_ch > p2.obj[len(p2.obj)-1] {
						vasPoste(p2, r)
					} else {
						dch_obj(p3, r)
						vasPoste(p2, r)
					}
				}
			}
		default:
			if r.flag_ch == true {
				dch_obj(p3, r)
				if f2 == 0 {
					ch_obj(p3, r)
				} else {
					if p3.obj[0] > p2.obj[len(p2.obj)-1] {
						ch_obj(p3, r)
					}
				}
				vasPoste(p2, r)
			} else {
				if f2 == 0 {
					ch_obj(p3, r)
				} else {
					if p3.obj[0] > p2.obj[len(p2.obj)-1] {
						ch_obj(p3, r)
					}
				}
				vasPoste(p2, r)
			}
		}
		mutex2.Unlock()
	}
	if stop(s, p2) == false {
		logique(p1, p2, p3, r, s)
	}
}

func logique(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, s []int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Logique failed:", err)
			os.Exit(300)
		}
	}()
	if stop(s, p2) == false {
		switch r.tip {
		case 13:
			switch r.nr_poste {
			case 1:
				logique_1(p1, p2, p3, r, s)
			case 3:
				logique_3(p1, p2, p3, r, s)
			}
		case 32:
			switch r.nr_poste {
			case 2:
				logique_2(p1, p2, p3, r, s)
			case 3:
				logique_3(p1, p2, p3, r, s)
			}
		}
	}
}

func stop(s []int, p2 *Poste) bool {
	if len(s) == len(p2.obj) {
		result := "DONE!!! --->"
		//fmt.Print("DONE!!! ---> ")
		for i := 0; i < len(p2.obj); i++ {
			//fmt.Printf("%d ", p2.obj[i])
			result += " " + strconv.Itoa(p2.obj[i])
		}
		//fmt.Printf("\n")
		terminoiu <- result
		return true
	}
	return false
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(8)
}

func main() {

	//Random
	// k:= []int{1,2,3,4,5,6,7,8,9}
	// #pieces au poste 1 <- rand
	// #pieces au poste 3 <- rand
	// #pieces au poste 1 + #pieces au poste 3 <= len(k)
	// #pieces au poste 2 <- [0,1] autrement on peut avoir des erreurs
	// pour chaque poste
	//  index = rand(0, len(k))
	//  s = k[index]
	//  k = k[:index-1] + k[index:]

	k = []int{1, 2, 3, 4, 5, 6, 7, 8}

	NO_P_P2 := rand.Intn(2)
	NO_P_P1 := rand.Intn(len(k) - NO_P_P2)
	NO_P_P3 := len(k) - NO_P_P1 - NO_P_P2

	//fmt.Printf("#1: %d \n#2: %d \n#3: %d\n", NO_P_P1, NO_P_P2, NO_P_P3)

	s := make([]int, NO_P_P1)
	s1 := make([]int, NO_P_P2)
	s2 := make([]int, NO_P_P3)
	for i := 0; i < NO_P_P1; i = i + 1 {
		index := rand.Intn(len(k))
		s[i] = k[index]
		k[index], k = k[len(k)-1], k[:len(k)-1]
	}
	for i := 0; i < NO_P_P2; i = i + 1 {
		index := rand.Intn(len(k))
		s1[i] = k[index]
		k[index], k = k[len(k)-1], k[:len(k)-1]
	}
	for i := 0; i < NO_P_P3; i = i + 1 {
		index := rand.Intn(len(k))
		s2[i] = k[index]
		k[index], k = k[len(k)-1], k[:len(k)-1]
	}

	P1 := Poste{1, s}
	P2 := Poste{2, s1}
	P3 := Poste{3, s2}
	s = append(s, s1...)
	s = append(s, s2...)

	r1 := Robot{1, 30000, false, 30000, 1, 13}
	r2 := Robot{2, 30000, false, 30000, 3, 32}
	r3 := Robot{3, 30000, false, 30000, 1, 13}
	r4 := Robot{4, 30000, false, 30000, 3, 32}
	r5 := Robot{5, 30000, false, 30000, 1, 13}
	r6 := Robot{6, 30000, false, 30000, 3, 32}
	r7 := Robot{7, 30000, false, 30000, 1, 13}
	r8 := Robot{8, 30000, false, 30000, 3, 32}

	go logique(&P1, &P2, &P3, &r1, s)
	go logique(&P1, &P2, &P3, &r2, s)
	go logique(&P1, &P2, &P3, &r3, s)
	go logique(&P1, &P2, &P3, &r4, s)
	go logique(&P1, &P2, &P3, &r5, s)
	go logique(&P1, &P2, &P3, &r6, s)
	go logique(&P1, &P2, &P3, &r7, s)
	go logique(&P1, &P2, &P3, &r8, s)

	fmt.Println(<-terminoiu)
	//<-terminoiu
}
