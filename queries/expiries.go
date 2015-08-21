package queries

// Turn Number Query

func (q *turnNumberQuery) isExpired(now interface{}) bool {
	return false
}

func (q *turnNumberQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Board State at Turn Query

func (q *boardStateAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *boardStateAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Move at Turn Query

func (q *moveAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *moveAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Draw Offer State Query

func (q *drawOfferStateQuery) isExpired(now interface{}) bool {
	return false
}

func (q *drawOfferStateQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// User Games Query

func (q *userGamesQuery) isExpired(now interface{}) bool {
	return false
}

func (q *userGamesQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Game Players Query

func (q *gamePlayersQuery) isExpired(now interface{}) bool {
	return false
}

func (q *gamePlayersQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Valid Moves At Turn Query

func (q *validMovesAtTurnQuery) isExpired(now interface{}) bool {
	return false
}

func (q *validMovesAtTurnQuery) getExpiration(now interface{}) interface{} {
	return nil
}

// Game Query

func (q *gameQuery) isExpired(now interface{}) bool {
	return false
}

func (q *gameQuery) getExpiration(now interface{}) interface{} {
	return nil
}
