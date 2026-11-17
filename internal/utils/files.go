package utils

import "path/filepath"

// ProjectFilePath возвращает абсолютный путь к файлу относительно корня проекта (папки с go.mod)
func ProjectFilePath(relativePath string) string {
	return filepath.Join(findGoModDir(searchCallerFile()), relativePath)
}
