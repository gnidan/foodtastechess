package game

type Piece interface{
	Name() string
	Color() Color
	Moves() []Move
}

type Color int
const(
	White Color = iota
	Black
)

type Pawn struct {
	color Color
	moves []Move
}
func (p *Pawn) Color{ return p.color }
func (p *Pawn) Moves(){ return p.moves }

type Rook struct {
	color Color
	moves []Move
}
func (r *Rook) Color{ return r.color }
func (r *Rook) Moves(){ return r.moves }

type Knight struct {
	color Color
	moves []Move
}
func (n *Knight) Color{ return n.color }
func (n *Knight) Moves(){ return n.moves }

type Bishop struct {
	color Color
	moves []Move
}
func (b *Bishop) Color{ return b.color }
func (b *Bishop) Moves(){ return b.moves }

type Queen struct {
	color Color
	moves []Move
}
func (q *Queen) Color{ return q.color }
func (q *Queen) Moves(){ return q.moves }

type King struct {
	color Color
	moves []Move
}
func (k *King) Color{ return k.color }
func (k *King) Moves(){ return k.moves }



func NewPawn(color Color) Pawn{
	p := Pawn{}
	p.color = color
	p.moves := [...]{
		new AdvancingMove(1, /*TODO: based on color */), //move forward
		new CapturingMove(-1, 1), new CapturingMove(1, 1), //Diagonal capturing
		}
	if(p.color == White){
		p.moves.append( 
			new EnPassantMove(-1, 1), new EnPassantMove(1, 1)
		)
	}else{
		p.moves.append( 
			new EnPassantMove(-1, -1), new EnPassantMove(1, -1)
		)
	}
	return p
}//NewPawn

func NewRook(color Color) Rook{
	r := Rook{}
	r.color = color
	r.moves := [...]{
		new UnboundMove(1, 0), new UnboundMove(0, -1), //Horizontals
		new UnboundMove(0, 1), new UnboundMove(0, -1), //Verticals
		new CastlingMove( /*TODO*/)
		}
	return r
}//NewRook

func NewKnight(color Color) Knight{
	n := Knight{}
	n.color = color
	n.moves := [...]{
		new BoundMove(-2, 1), new BoundMove(-1, 2), //up-left
		new BoundMove(1, 2), new BoundMove(2, 1), //up-right
		new BoundMove(-2, 1), new BoundMove(-1, 2), //down-left
		new BoundMove(-2, 1), new BoundMove(-1, 2) //down-right
		}
	return n
}//NewKnight

func NewBishop(color Color) Bishop{
	b := Bishop{}
	b.color = color
	b.moves := [...]{
		new UnboundMove(-1, 1), new UnboundMove(1, 1), //up diagonals
		new UnboundMove(-1, -1), new UnboundMove(1, -1) //down diagonals
		}
	return b
}//NewBishop

func NewQueen(color Color) Queen{
	q := Queen{}
	q.color = color
	q.moves := [...]{
		new UnboundMove(-1, 1), new UnboundMove(0, 1), new UnboundMove(1, 1), //forward
		new UnboundMove(-1, 0), new UnboundMove(1, 0), //sideways
		new UnboundMove(-1, -1), new UnboundMove(0, -1), new UnboundMove(1, -1) //backwards
		}
	return q
}//NewQueen

func NewKing(color Color) King{
	k := King{}
	k.color = color
	k.moves := [...]{
		new BoundMove(-1, 1), new BoundMove(0, 1), new BoundMove(1, 1), //forward
		new BoundMove(-1, 0), new BoundMove(1, 0), //sideways
		new BoundMove(-1, -1), new BoundMove(0, -1), new BoundMove(1, -1) //backwards
		}
	return k
}//NewKing
