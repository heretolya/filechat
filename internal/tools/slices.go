package tools

import "errors"

func SliceToChunks[T any](
	slice []T,
	chunkSize,
	overlap int,
) ([][]T, error) {
	if chunkSize == 0 {
		errMsg := "chunk size is zero"
		return nil, errors.New(errMsg)
	}
	if overlap == chunkSize {
		errMsg := "overlap >= chunk size"
		return nil, errors.New(errMsg)
	}
	chunks := make([][]T, 0)
	if chunkSize >= len(slice) {
		chunks = append(chunks, slice)
		return chunks, nil
	}
	step := chunkSize - overlap
	for i := 0; i < len(slice); i += step {
		end := min(i+chunkSize, len(slice))
		chunks = append(chunks, slice[i:end])
		if end == len(slice) {
			break
		}
	}
	return chunks, nil
}
