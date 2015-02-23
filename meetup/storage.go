package meetup

import "github.com/jmcvetta/neoism"

var people, techs map[string]int

func (a API) storeInDB(m Member) (err error) {
	if people == nil || techs == nil {
		people = make(map[string]int)
		techs = make(map[string]int)
	}

	var mn, tn *neoism.Node

	if id, ok := people[m.Name]; ok {
		mn, err = a.db.Node(id)
	} else {
		mn, err = a.db.CreateNode(neoism.Props{"name": m.Name})
	}
	if err != nil {
		return err
	}

	for _, t := range m.Topics {
		if id, ok := techs[t.Name]; ok {
			tn, err = a.db.Node(id)
		} else {
			tn, err = a.db.CreateNode(neoism.Props{"name": t.Name})
		}
		if err != nil {
			return err
		}

		mn.Relate("interested_in", tn.Id(), nil)
	}

	return err
}
