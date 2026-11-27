package concept

func Init(path, slug string) (string, error) {
	return Layout.Create(path, slug, Template)
}
