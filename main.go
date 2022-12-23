package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const dataFile = "pwm.data"

type Entry struct {
	Username string
	Password string
	URL      string
	Notes    string
}

func main() {
	var add bool
	flag.BoolVar(&add, "add", false, "Añade una nueva entrada al manejador de contraseñas")

	var list bool
	flag.BoolVar(&list, "list", false, "Lista todas las entradas del manejador de contraseñas")

	flag.Parse()

	if add {
		entry := readEntry()
		saveEntry(entry)
		fmt.Println("Entrada guardada con éxito")
	} else if list {
		entries := loadEntries()
		if entries == nil {
			fmt.Println("No hay entradas guardadas")
		} else {
			for _, entry := range entries {
				fmt.Println("Username:", entry.Username)
				fmt.Println("Password:", entry.Password)
				fmt.Println("URL:", entry.URL)
				fmt.Println("Notes:", entry.Notes)
				fmt.Println()
			}
		}
	} else {
		fmt.Println("Uso: pwm [opciones] ")
		fmt.Println("Opciones:")
		fmt.Println("  -add\tAñade una nueva entrada al manejador de contraseñas")
		fmt.Println("  -list\tLista todas las entradas del manejador de contraseñas")
	}
}

func readEntry() *Entry {
	entry := &Entry{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	entry.Username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	entry.Password = strings.TrimSpace(password)

	fmt.Print("URL: ")
	url, _ := reader.ReadString('\n')
	entry.URL = strings.TrimSpace(url)

	fmt.Print("Notes: ")
	notes, _ := reader.ReadString('\n')
	entry.Notes = strings.TrimSpace(notes)

	return entry
}

func saveEntry(entry *Entry) {
	file, err := os.OpenFile(dataFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		fmt.Println("Error al abrir el archivo", err)
		return
	}
	defer file.Close()

	type Data struct {
		Entries []*Entry
	}
	data := &Data{Entries: []*Entry{entry}}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al convertir a JSON", err)
		return
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println("Error al obtener información del archivo", err)
		return
	}

	if info.Size() == 0 {
		_, err = file.Write(jsonData)
	} else {
		_, err = file.Write(append([]byte(","), jsonData...))
	}

	if err != nil {
		fmt.Println("Error al escribir en el archivo", err)
		return
	}
}

func loadEntries() []*Entry {
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Error al abrir el archivo", err)
		return nil
	}
	defer file.Close()

	type Data struct {
		Entries []*Entry
	}
	data := &Data{}

	jsonData, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Error al leer el archivo", err)
		return nil
	}

	err = json.Unmarshal(jsonData, data)
	if err != nil {
		fmt.Println("Error al convertir de JSON", err)
		return nil
	}

	return data.Entries
}