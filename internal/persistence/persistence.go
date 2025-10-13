package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ParkhomenkoDV/URLShortener/internal/model"
)

// Cериализация/десериализация данных в JSON
type JSONPersistence interface {
	Save(filePath string, data map[string]string) error
	Load(filePath string) (map[string]string, int, error)
}

// Реализация для работы с JSON файлами
type FileJSONPersistence struct{}

// Создаёт новый экземпляр FileJSONPersistence
func NewFileJSONPersistence() *FileJSONPersistence {
	return &FileJSONPersistence{}
}

// Сохраняет данные в JSON файл
func (p *FileJSONPersistence) Save(filePath string, data map[string]string) error {
	var records []model.URLRecord

	counter := 1
	for shortURL, originalURL := range data {
		record := model.URLRecord{
			ID:          counter,
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		}
		records = append(records, record)
		counter++
	}

	return p.saveRecordsToFile(filePath, records)
}

// Загружает данные из JSON файла
func (p *FileJSONPersistence) Load(filePath string) (map[string]string, int, error) {
	records, err := p.loadRecordsFromFile(filePath)
	if err != nil {
		return nil, 0, err
	}

	data := make(map[string]string)
	maxID := 0
	for _, record := range records {
		data[record.ShortURL] = record.OriginalURL
		if record.ID > maxID {
			maxID = record.ID
		}
	}

	return data, maxID, nil
}

// Сохраняет записи в файл
func (p *FileJSONPersistence) saveRecordsToFile(filePath string, records []model.URLRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Создаем директорию если её нет
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// Загружает записи из файла
func (p *FileJSONPersistence) loadRecordsFromFile(filePath string) ([]model.URLRecord, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []model.URLRecord{}, nil // Если файл не существует, возвращаем пустой массив
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Пустой файл
	if len(data) == 0 {
		return []model.URLRecord{}, nil
	}

	var records []model.URLRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}
