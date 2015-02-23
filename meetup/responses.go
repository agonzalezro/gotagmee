package meetup

type GroupsResponse struct {
	Results []struct {
		Members int
	}
}

type MembersResponse struct {
	Results []struct {
		Name   string
		Topics []struct {
			ID           int
			Name, URLKey string
		}
	}
}
