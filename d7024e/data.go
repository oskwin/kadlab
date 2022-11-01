package d7024e

type DataStorage struct {
	Data map[string]*File
}

type File struct {
	Key   *string
	Value []byte
}
