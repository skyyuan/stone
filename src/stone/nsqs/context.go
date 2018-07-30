package nsqs

type ctxMap map[string]interface{}
type context struct {
	store ctxMap
}

func newContext() *context {
	return &context{
		store: make(ctxMap),
	}
}

func (c *context) Get(key string) interface{} {
	return c.store[key]
}

func (c *context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(ctxMap)
	}
	c.store[key] = val
}
