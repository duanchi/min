package library

type influxResult struct {
	// One entry in both slices is created for every executed statement result.
	affectedRows int64
	insertIds    []int64
}

func (res *influxResult) LastInsertId() (int64, error) {
	return res.insertIds[len(res.insertIds)-1], nil
}

func (res *influxResult) RowsAffected() (int64, error) {
	return res.affectedRows, nil
}
