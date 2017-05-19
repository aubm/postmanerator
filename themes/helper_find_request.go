package themes

import "github.com/aubm/postmanerator/postman"

func helperFindRequest(requests []postman.Request, ID string) *postman.Request {
	for _, r := range requests {
		if r.ID == ID {
			return &r
		}
	}
	return nil
}
