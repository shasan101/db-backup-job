package uploaders

import "context"

type LocalFs struct {
}

func InitLocal(c context.Context, test1, test2 string) (*LocalFs, error) {
	return &LocalFs{}, nil
}

func (l *LocalFs) UploadBackup(c context.Context, filePath string) error {
	return nil
}
