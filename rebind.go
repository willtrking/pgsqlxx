package pgsqlxx

import "strconv"

// From https://github.com/jmoiron/sqlx/blob/398dd5876282499cdfd4cb8ea0f31a672abe9495/bind.go#L40
func Rebind(query string) string {

	qb := []byte(query)
	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(qb)+10)
	j := 1
	for _, b := range qb {
		if b == '?' {
			rqb = append(rqb, '$')
			for _, b := range strconv.Itoa(j) {
				rqb = append(rqb, byte(b))
			}
			j++
		} else {
			rqb = append(rqb, b)
		}
	}

	return string(rqb)
}
