package game

type Piece interface {
	Name() string
	Color() Color
	Moves() []Move
}

type Color int

const (
	White Color = iota
	Black
)

type Pawn struct {
	color Color
	moves []Move
}

func (p *Pawn) Color() Color  { return p.color }
func (p *Pawn) Moves() []Move { return p.moves }

type Rook struct {
	color Color
	moves []Move
}

func (r *Rook) Color() Color  { return r.color }
func (r *Rook) Moves() []Move { return r.moves }

type Knight struct {
	color Color
	moves []Move
}

func (n *Knight) Color() Color  { return n.color }
func (n *Knight) Moves() []Move { return n.moves }

type Bishop struct {
	color Color
	moves []Move
}

func (b *Bishop) Color() Color  { return b.color }
func (b *Bishop) Moves() []Move { return b.moves }

type Queen struct {
	color Color
	moves []Move
}

func (q *Queen) Color() Color  { return q.color }
func (q *Queen) Moves() []Move { return q.moves }

type King struct {
	color Color
	moves []Move
}

func (k *King) Color() Color  { return k.color }
func (k *King) Moves() []Move { return k.moves }

func NewPawn(color Color) Pawn {
	p := Pawn{}
	p.color = color
	p.moves = []Move{
		&AdvancingMove{1 /*TODO: based on color */}, //move forward
		&CapturingMove{1, 1}, &CapturingMove{1, 1},  //Diagonal capturing
	}
	if p.color == White {
		p.moves = append(p.moves,
			&EnPassantMove{-1, 1}, &EnPassantMove{1, 1},
		)
	} else {
		p.moves = append(p.moves,
			&EnPassantMove{-1, -1}, &EnPassantMove{1, -1},
		)
	}
	return p
} //NewPawn

func NewRook(color Color) Rook {
	r := Rook{}
	r.color = color
	r.moves = []Move{
		&UnboundMove{1, 0}, &UnboundMove{0, -1}, //Horizontals
		&UnboundMove{0, 1}, &UnboundMove{0, -1}, //Verticals
		&CastlingMove{ /*TODO*/ },
	}
	return r
} //NewRook

func NewKnight(color Color) Knight {
	n := Knight{}
	n.color = color
	n.moves = []Move{
		&BoundMove{-2, 1}, &BoundMove{-1, 2}, //up-left
		&BoundMove{1, 2}, &BoundMove{2, 1}, //up-right
		&BoundMove{-2, 1}, &BoundMove{-1, 2}, //down-left
		&BoundMove{-2, 1}, &BoundMove{-1, 2}, //down-right
	}
	return n
} //NewKnight

func NewBishop(color Color) Bishop {
	b := Bishop{}
	b.color = color
	b.moves = []Move{
		&UnboundMove{-1, 1}, &UnboundMove{1, 1}, //up diagonals
		&UnboundMove{-1, -1}, &UnboundMove{1, -1}, //down diagonals
	}
	return b
} //NewBishop

func NewQueen(color Color) Queen {
	q := Queen{}
	q.color = color
	q.moves = []Move{
		&UnboundMove{-1, 1}, &UnboundMove{0, 1}, &UnboundMove{1, 1}, //forward
		&UnboundMove{-1, 0}, &UnboundMove{1, 0}, //sideways
		&UnboundMove{-1, -1}, &UnboundMove{0, -1}, &UnboundMove{1, -1}, //backwards
	}
	return q
} //NewQueen

func NewKing(color Color) King {
	k := King{}
	k.color = color
	k.moves = []Move{
		&BoundMove{-1, 1}, &BoundMove{0, 1}, &BoundMove{1, 1}, //forward
		&BoundMove{-1, 0}, &BoundMove{1, 0}, //sideways
		&BoundMove{-1, -1}, &BoundMove{0, -1}, &BoundMove{1, -1}, //backwards
	}
	return k
} //NewKing
