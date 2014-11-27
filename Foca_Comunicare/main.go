// Foca_Comunicare project main.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Robot struct {
	id       int
	obj_ch   int
	flag_ch  bool
	mp2      int
	nr_poste int
}

type Poste struct {
	nr  int
	obj []int
}

type Info struct {
	id          int
	post        int
	nr_incarcat int
}

var ch chan *Info

func ch_obj(p *Poste, r *Robot, i *Info) {
	r.flag_ch = true
	r.obj_ch = p.obj[0]
	if len(p.obj) > 1 {
		p.obj = p.obj[1:]
	} else {
		p.obj = []int{}
	}
	if r.mp2 == 30000 {
		r.mp2 = r.obj_ch
	}
	fmt.Printf("Incarcare robot #%d cu : %d de la postul %d %v\n", r.id, r.obj_ch, p.nr, i)
	if len(p.obj) > 0 {
		fmt.Printf("Postul %d contine: ", p.nr)
		for i := 0; i < len(p.obj); i++ {
			fmt.Printf("%d ", p.obj[i])
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Postul %d este gol %v\n", p.nr, i)
	}
}

func ch_obj_fin(p *Poste, r *Robot, i *Info) {
	r.flag_ch = true
	f := len(p.obj) - 1
	r.obj_ch = p.obj[f]
	if len(p.obj) > 1 {
		p.obj = p.obj[0:f]
	} else {
		p.obj = []int{}
	}
	fmt.Printf("Incarcare robot #%d cu : %d de la postul %d %v\n", r.id, r.obj_ch, p.nr, i)
	if len(p.obj) > 0 {
		fmt.Printf("Postul %d contine: ", p.nr)
		for i := 0; i < len(p.obj); i++ {
			fmt.Printf("%d ", p.obj[i])
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Postul %d este gol %v\n", p.nr, i)
	}
}

func dch_obj(p *Poste, r *Robot, i *Info) {
	p.obj = append(p.obj, r.obj_ch)
	fmt.Printf("Robotul #%d descarca : %d la postul %d %v\n", r.id, r.obj_ch, p.nr, i)
	if len(p.obj) > 0 {
		fmt.Printf("Postul %d contine: ", p.nr)
		for i := 0; i < len(p.obj); i++ {
			fmt.Printf("%d ", p.obj[i])
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Postul %d este gol %v\n", p.nr, i)
	}
	r.flag_ch = false
	r.obj_ch = 30000
}

func vasPoste(p *Poste, r *Robot, i *Info) {
	x := rand.Float32()
	switch p.nr {
	case 1:
		switch {
		case x < .10:
			r.nr_poste = 2
		case x < .20:
			r.nr_poste = 3
		default:
			r.nr_poste = 1
		}
	case 2:
		switch {
		case x < .10:
			r.nr_poste = 1
		case x < .20:
			r.nr_poste = 3
		default:
			r.nr_poste = 2
		}
	case 3:
		switch {
		case x < .10:
			r.nr_poste = 1
		case x < .20:
			r.nr_poste = 2
		default:
			r.nr_poste = 3
		}
	}
	fmt.Printf("Robotul #%d merge la postul : %d %f %v\n", r.id, r.nr_poste, x, i)
}

func logique_1(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, i *Info, s []int) {
	f := 0
	if len(p2.obj) > 0 {
		f = len(p2.obj) - 1
	}
	if r.flag_ch == true {
		if f == 0 {
			vasPoste(p2, r, i)
			if r.nr_poste == 2 {
				verificaSiTrimite(r, i)
			}
			logique(p1, p2, p3, r, i, s)
		} else if f != 0 {
			if r.obj_ch > p2.obj[f] {
				vasPoste(p2, r, i)
				if r.nr_poste == 2 {
					verificaSiTrimite(r, i)
				}
				logique(p1, p2, p3, r, i, s)
			}
		} else {
			vasPoste(p3, r, i)
			logique(p1, p2, p3, r, i, s)
		}
	} else {
		if len(p1.obj) != 0 {
			ch_obj(p1, r, i)
			if r.obj_ch < r.mp2 {
				r.mp2 = r.obj_ch
				vasPoste(p3, r, i)
				logique(p1, p2, p3, r, i, s)
			} else if f != 0 && r.obj_ch < p2.obj[f] {
				vasPoste(p3, r, i)
				logique(p1, p2, p3, r, i, s)
			} else {
				vasPoste(p2, r, i)
				if r.nr_poste == 2 {
					verificaSiTrimite(r, i)
				}
				logique(p1, p2, p3, r, i, s)
			}
		} else {
			vasPoste(p3, r, i)
			logique(p1, p2, p3, r, i, s)
		}
	}
}

func logique_2(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, i *Info, s []int) {
	f := len(p2.obj)
	if r.flag_ch == true {
		switch {
		case f != 0:
			{
				if r.obj_ch > p2.obj[len(p2.obj)-1] {
					//var info *Info
					info1 := new(Info)
					info1.id = r.id
					info1.post = -1
					info1.nr_incarcat = 30000
					//info = Info{r.id, -1, 30000}
					ch <- info1
					dch_obj(p2, r, i)

					if stop(s, p2) == false {
						logique(p1, p2, p3, r, i, s)
					}
				}
			}
		case f == 0:
			{
				//var info *Info
				info1 := new(Info)
				info1.id = r.id
				info1.post = -1
				info1.nr_incarcat = 30000
				//info = Info{r.id, -1, 30000}
				ch <- info1
				dch_obj(p2, r, i)

				if stop(s, p2) == false {
					logique(p1, p2, p3, r, i, s)
				}
			}
		}
	}
	if f != 0 {
		f = len(p2.obj) - 1
		if len(p3.obj) == 0 {
			vasPoste(p1, r, i)
			logique(p1, p2, p3, r, i, s)
		} else if p2.obj[f] > p3.obj[0] && r.flag_ch == false {
			ch_obj_fin(p2, r, i)
			if r.obj_ch < r.mp2 {
				r.mp2 = r.obj_ch
			}
			vasPoste(p3, r, i)
			logique(p1, p2, p3, r, i, s)
		} else {
			vasPoste(p3, r, i)
			logique(p1, p2, p3, r, i, s)
		}
	}
}

func logique_3(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, i *Info, s []int) {
	if r.flag_ch == true {
		dch_obj(p3, r, i)
	}
	if len(p2.obj) == 0 {
		ch_obj(p3, r, i)
		if r.obj_ch < r.mp2 {
			r.mp2 = r.obj_ch
		}
		vasPoste(p2, r, i)
		if r.nr_poste == 2 {
			verificaSiTrimite(r, i)
		}
		logique(p1, p2, p3, r, i, s)
	} else if len(p3.obj) > 0 {
		if p3.obj[0] > p2.obj[len(p2.obj)-1] {
			ch_obj(p3, r, i)
			if r.obj_ch < r.mp2 {
				r.mp2 = r.obj_ch
			}
			vasPoste(p2, r, i)
			if r.nr_poste == 2 {
				verificaSiTrimite(r, i)
			}
			logique(p1, p2, p3, r, i, s)
		} else {
			vasPoste(p2, r, i)
			logique(p1, p2, p3, r, i, s)
		}
	} else {
		//vasPoste(p2, r)
		logique(p1, p2, p3, r, i, s)
	}
}

func logique(p1 *Poste, p2 *Poste, p3 *Poste, r *Robot, i *Info, s []int) {
	if stop(s, p2) == false {
		switch r.nr_poste {
		case 1:
			logique_1(p1, p2, p3, r, i, s)
		case 2:
			logique_2(p1, p2, p3, r, i, s)
		case 3:
			logique_3(p1, p2, p3, r, i, s)
		}
	}
}

func stop(s []int, p2 *Poste) bool {
	if len(s) == len(p2.obj) {
		return true
	}
	return false
}

func stop_robot(r *Robot) {
	fmt.Println("Robotul #%d sta pe loc\n", r.id)
}

func verificaSiTrimite(r *Robot, i *Info) {
	if r.id != i.id && r.nr_poste == i.post && r.obj_ch > i.nr_incarcat {
		stop_robot(r)
		/*for r.id != i.id && r.nr_poste == i.post && r.obj_ch > i.nr_incarcat {
		}*/
	}
	//var info1 *Info
	info1 := new(Info)
	info1.id = r.id
	info1.post = r.nr_poste
	info1.nr_incarcat = r.obj_ch
	ch <- info1
}

func deamon(i *Info) {
	for {
		i1 := new(Info)
		i1 = <-ch
		i.id = i1.id
		i.post = i1.post
		i.nr_incarcat = i1.nr_incarcat
		fmt.Println(i)
	}
}

func main() {
	s := []int{4, 3, 1, 2}
	s1 := []int{7, 6}
	s2 := []int{5, 8}
	P1 := Poste{1, s}
	P2 := Poste{2, s1}
	P3 := Poste{3, s2}
	s = append(s, s1...)
	s = append(s, s2...)
	var i Info
	i.id = -1
	i.post = -1
	i.nr_incarcat = 30000
	ch = make(chan *Info)
	r1 := Robot{1, 30000, false, 30000, 1}
	r2 := Robot{2, 30000, false, 30000, 1}
	//r3 := Robot{3, 30000, false, 30000, 1}
	//r4 := Robot{4, 30000, false, 30000, 1}
	go deamon(&i)
	go logique(&P1, &P2, &P3, &r2, &i, s)
	go logique(&P1, &P2, &P3, &r1, &i, s)
	//go logique(&P1, &P2, &P3, &r3, &i, s)
	//go logique(&P1, &P2, &P3, &r4, &i, s)
	time.Sleep(10 * time.Millisecond)
}
