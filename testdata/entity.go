package entity

import "time"

type User struct {
	ID        uint64     `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	NonTag    uint
}

type Value uint32
type StrongItem struct {
	ID    uint64  `db:"id"`
	Name  *string `db:"name"`
	Value Value   `db:"value"`
}

type NonDBStruct struct {
	Power uint64
}

func f1() {
	type User [1]int
}

func f2() {
	type User struct{ k int }
	_ = User{k: 0}
}
