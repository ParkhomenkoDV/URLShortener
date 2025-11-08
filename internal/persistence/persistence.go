package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ParkhomenkoDV/URLShortener/internal/model"
)

// JSONPersistence определяет интерфейс для сериализации и десериализации данных в формате JSON.
// Используется для абстракции над различными реализациями хранения данных (файлы, базы данных и т.д.).
type JSONPersistence interface {
	// Save сохраняет данные URL и соответствующие им пользовательские идентификаторы в хранилище.
	// Параметры:
	//   - filePath: путь к файлу для сохранения
	//   - data: маппинг коротких URL на оригинальные URL
	//   - userMap: маппинг коротких URL на идентификаторы пользователей
	// Возвращает ошибку если операция сохранения не удалась.
	Save(filePath string, data map[string]string, userMap map[string]string) error

	// Load загружает данные URL и пользовательские идентификаторы из хранилища.
	// Параметры:
	//   - filePath: путь к файлу для загрузки
	// Возвращает:
	//   - data: маппинг коротких URL на оригинальные URL
	//   - userMap: маппинг коротких URL на идентификаторы пользователей
	//   - maxID: максимальный ID записи (для продолжения нумерации)
	//   - error: ошибка если операция загрузки не удалась
	Load(filePath string) (map[string]string, map[string]string, int, error)
}

// FileJSONPersistence реализует JSONPersistence интерфейс для работы с JSON файлами.
// Сохраняет данные в виде массива URLRecord объектов в формате JSON.
type FileJSONPersistence struct{}

// NewFileJSONPersistence создаёт и возвращает новый экземпляр FileJSONPersistence.
// Используется для инъекции зависимости в компоненты, требующие персистентности.
func NewFileJSONPersistence() *FileJSONPersistence {
	return &FileJSONPersistence{}
}

// Save сохраняет данные URL и пользовательские идентификаторы в JSON файл.
// Преобразует переданные мапы в массив URLRecord структур и сохраняет их в файл.
// Автоматически создает директории если они не существуют.
//
// Пример сохраняемого JSON:
//
//	[
//	  {
//	    "id": 1,
//	    "short_url": "abc123",
//	    "original_url": "https://example.com",
//	    "user_id": "user-1"
//	  }
//	]
func (p *FileJSONPersistence) Save(filePath string, data map[string]string, userMap map[string]string) error {
	var records []model.URLRecord

	// Конвертируем мапы в массив URLRecord структур с последовательной нумерацией ID
	counter := 1
	for shortURL, originalURL := range data {
		userID := userMap[shortURL]
		record := model.URLRecord{
			ID:          counter,
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
		}
		records = append(records, record)
		counter++
	}

	return p.saveRecordsToFile(filePath, records)
}

// Load загружает данные URL и пользовательские идентификаторы из JSON файла.
// Если файл не существует, возвращает пустые мапы и maxID = 0.
// Восстанавливает маппинги из массива URLRecord структур и вычисляет максимальный ID.
func (p *FileJSONPersistence) Load(filePath string) (map[string]string, map[string]string, int, error) {
	records, err := p.loadRecordsFromFile(filePath)
	if err != nil {
		return nil, nil, 0, err
	}

	// Восстанавливаем маппинги из массива записей
	data := make(map[string]string)
	userMap := make(map[string]string)
	maxID := 0
	for _, record := range records {
		data[record.ShortURL] = record.OriginalURL
		userMap[record.ShortURL] = record.UserID
		if record.ID > maxID {
			maxID = record.ID
		}
	}

	return data, userMap, maxID, nil
}

// saveRecordsToFile сохраняет массив URLRecord в JSON файл с отступами для читаемости.
// Автоматически создает необходимые директории. Если records nil, сохраняет пустой массив.
//
// Внутренний метод, используется реализацией Save.
func (p *FileJSONPersistence) saveRecordsToFile(filePath string, records []model.URLRecord) error {
	// Обрабатываем случай nil records - сохраняем как пустой массив
	if records == nil {
		records = []model.URLRecord{}
	}

	// Сериализуем в JSON с отступами для удобства чтения и отладки
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Создаем все необходимые директории рекурсивно
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Записываем данные в файл с правами чтения/записи для владельца
	return os.WriteFile(filePath, data, 0644)
}

// loadRecordsFromFile загружает массив URLRecord из JSON файла.
// Если файл не существует, возвращает пустой массив без ошибки.
// Если файл пустой, возвращает пустой массив без ошибки.
//
// Внутренний метод, используется реализацией Load.
func (p *FileJSONPersistence) loadRecordsFromFile(filePath string) ([]model.URLRecord, error) {
	// Проверяем существование файла - если не существует, возвращаем пустой массив
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []model.URLRecord{}, nil
	}

	// Читаем содержимое файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Обрабатываем случай пустого файла
	if len(data) == 0 {
		return []model.URLRecord{}, nil
	}

	// Десериализуем JSON в массив URLRecord
	var records []model.URLRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	return records, nil
}
