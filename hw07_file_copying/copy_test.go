package main

import (
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем данные во временный файл
	_, err = tmpFile.WriteString("Hello, World!")
	if err != nil {
		t.Fatal(err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Копируем файл
	err = Copy(tmpFile.Name(), "test_copy.txt", 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что файл скопирован правильно
	copyFile, err := os.Open("test_copy.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer copyFile.Close()

	copyFileInfo, err := copyFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if copyFileInfo.Size() != 13 {
		t.Errorf("Expected file size to be 13, but got %d", copyFileInfo.Size())
	}

	// Close the file before removing it
	copyFile.Close()

	// Удаляем созданный файл
	err = os.Remove("test_copy.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCopyWithOffset(t *testing.T) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем данные во временный файл
	_, err = tmpFile.WriteString("Hello, World!")
	if err != nil {
		t.Fatal(err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Копируем файл с отступом
	err = Copy(tmpFile.Name(), "test_copy.txt", 7, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что файл скопирован правильно
	copyFile, err := os.Open("test_copy.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer copyFile.Close()

	copyFileInfo, err := copyFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if copyFileInfo.Size() != 6 {
		t.Errorf("Expected file size to be 6, but got %d", copyFileInfo.Size())
	}

	// Close the file before removing it
	copyFile.Close()

	// Удаляем созданный файл
	err = os.Remove("test_copy.txt")
	if err != nil {
		t.Fatal(err)
	}
}
