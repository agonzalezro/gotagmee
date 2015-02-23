package db

import "github.com/jmcvetta/neoism"

type DB struct {
	conn *neoism.Database

	cachedTopics map[string]int
}

func NewDB(uri string) (*DB, error) {
	conn, err := neoism.Connect(uri)
	if err != nil {
		return nil, err
	}
	db := DB{conn: conn}
	db.cachedTopics = make(map[string]int)
	return &db, nil
}

func (db DB) Store(m Member) (err error) {
	var mn, tn *neoism.Node

	mn, err = db.conn.CreateNode(neoism.Props{"name": m.Name})
	mn.SetLabels([]string{"Member"})
	if err != nil {
		return err
	}

	for _, t := range m.Topics {
		if id, ok := db.cachedTopics[t]; ok {
			tn, err = db.conn.Node(id)
		} else {
			tn, err = db.conn.CreateNode(neoism.Props{"name": t})
			tn.SetLabels([]string{"Topic"})
			db.cachedTopics[t] = tn.Id()
		}
		if err != nil {
			return err
		}

		mn.Relate("interested_in", tn.Id(), nil)
	}

	return err
}
