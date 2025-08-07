package mentions

import (
	"regexp"

	"github.com/skamenetskiy/messages/internal/entity"
)

func Find(content string, links map[string]uint64) entity.Mentions {
	if links == nil {
		return nil
	}
	mentions := make(entity.Mentions, 0)
	if v := findAt(content, links); len(v) > 0 {
		mentions = append(mentions, v...)
	}
	if v := findHash(content, links); len(v) > 0 {
		mentions = append(mentions, v...)
	}
	if v := findDollar(content, links); len(v) > 0 {
		mentions = append(mentions, v...)
	}
	if v := findTilda(content, links); len(v) > 0 {
		mentions = append(mentions, v...)
	}
	return mentions
}

const (
	typeAt     uint32 = 1
	typeHash   uint32 = 2
	typeDollar uint32 = 3
	typeTilda  uint32 = 4
)

var (
	atR     = regexp.MustCompile("@([A-Za-z0-9_.]+)")
	hashR   = regexp.MustCompile("#([A-Za-z0-9_.]+)")
	dollarR = regexp.MustCompile("[$]([A-Za-z0-9_.]+)")
	tildaR  = regexp.MustCompile("~([A-Za-z0-9_.]+)")
)

func findAt(content string, links map[string]uint64) entity.Mentions {
	return find(content, links, typeAt, atR)
}

func findHash(content string, links map[string]uint64) entity.Mentions {
	return find(content, links, typeHash, hashR)
}

func findDollar(content string, links map[string]uint64) entity.Mentions {
	return find(content, links, typeDollar, dollarR)
}

func findTilda(content string, links map[string]uint64) entity.Mentions {
	return find(content, links, typeTilda, tildaR)
}

func find(content string, links map[string]uint64, t uint32, r *regexp.Regexp) entity.Mentions {
	res := r.FindAllStringIndex(content, -1)
	mentions := make(entity.Mentions, 0, len(res))
	for _, pos := range res {
		id, ok := links[content[pos[0]:pos[1]]]
		if !ok {
			continue
		}
		mentions = append(mentions, entity.Mention{
			ID:   id,
			Type: t,
			Pos: [2]uint32{
				uint32(pos[0]),
				uint32(pos[1]),
			},
		})
	}
	return mentions
}
