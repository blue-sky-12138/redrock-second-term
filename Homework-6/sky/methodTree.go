package sky

type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree

func (m methodTrees) getTree(method string) *node {
	for _, tree := range m {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}
