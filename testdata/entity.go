package entity

import "time"

type User struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	NonTag uint
}

func f1() {
	type User [1]int
}

func f2() {
	type User struct{ k int }
	_ = User{k: 0}
}
