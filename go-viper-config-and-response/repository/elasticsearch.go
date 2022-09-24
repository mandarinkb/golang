package repository

type ElasticSearchRepository interface {
	ESBuildJsonToByte(data interface{}) ([]byte, error)
	ESPost(index, indexType string, msg []byte) ([]byte, error)
}
