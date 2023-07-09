package clear

type Tokenizer interface {
	HasNext() bool
	Next() string
	Peek() string
}

type ssTzer struct {
	tokens []string
	index  int
}

func (ss *ssTzer) HasNext() bool {
	return ss.index < len(ss.tokens)
}

func (ss *ssTzer) Next() string {
	t := ss.tokens[ss.index]
	ss.index++
	return t
}

func (ss *ssTzer) Peek() string {
	return ss.tokens[ss.index]
}

type chainTzer struct {
	tzers []Tokenizer
}

func (ct *chainTzer) HasNext() bool {
	for _, t := range ct.tzers {
		if t.HasNext() {
			return true
		}
	}
	return false
}

func (ct *chainTzer) Next() string {
	for _, t := range ct.tzers {
		if t.HasNext() {
			return t.Next()
		}
	}
	return ""
}

func (ct *chainTzer) Peek() string {
	for _, t := range ct.tzers {
		if t.HasNext() {
			return t.Peek()
		}
	}
	return ""
}

func NewTokenizer(args ...string) Tokenizer {
	return &ssTzer{
		tokens: args,
	}
}

func NewTokenizerChain(tzers ...Tokenizer) Tokenizer {
	return &chainTzer{
		tzers: tzers,
	}
}
