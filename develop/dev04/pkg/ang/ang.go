package ang

import (
	"sort"
	"strings"
)

func byteSumFor(str string) int {
	b := []byte(str)
	output := 0
	for _, v := range b {
		output += int(v)
	}
	return output
}

func ShowAnagrams(in []string) map[string][]string {
	var output = make(map[string][]string)
mainLoop:
	for _, v := range in {
		v = strings.ToLower(v) // Все слова должны быть приведены к нижнему регистру.
		// Смотрим наличие нужного ключа. Для этого преобразуем строку в байты и посчитаем сумму этих байт.
		sum := byteSumFor(v)
		for k := range output {
			if byteSumFor(k) == sum && len(k) == len(v) { // Находим в словаре нужное значение
				output[k] = append(output[k], v)
				continue mainLoop
			}
		}
		output[v] = make([]string, 0, 1) // Если не нашли, добавляем ключ и инициализируем ему массив
	}
	for k := range output {
		if len(output[k]) == 0 {
			delete(output, k) // удаляем лишние ключи
			continue
		}
		sort.Strings(output[k])
	}
	return output
}
