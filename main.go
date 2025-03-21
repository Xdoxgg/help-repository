package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

type IdToId struct {
	Id1 int `json:"id1"`
	Id2 int `json:"id2"`
}

type Character struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Lore string `json:"lore"`
	Img  string `json:"img"`
}

type Team struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Orientation string      `json:"orientation"`
	Characters  []Character `json:"characters"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func connectDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return nil, err
	}
	// Проверяем соединение
	if err = db.Ping(); err != nil {
		fmt.Println("Не удалось подключиться к базе данных:", err)
		return nil, err
	}
	return db, nil
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Pages/_index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

//запросы в бд

func getAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	teams, err := getAllTeams(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
func getAllTeams(db *sql.DB) ([]Team, error) {
	namesQuery, err := db.Query("SELECT teams.id, teams.name, teams.orientation FROM teams")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer namesQuery.Close() // Закрываем запрос после использования

	var teams []Team
	for namesQuery.Next() {
		var team Team
		err := namesQuery.Scan(&team.ID, &team.Name, &team.Orientation)
		if err != nil {
			return nil, err
		}

		charactersQuery, err := db.Query("SELECT characters.id, characters.name, characters.role, characters.lore, characters.talents_build_emblems FROM teams JOIN team_to_character ON (teams.id = team_to_character.team_id) JOIN characters ON (characters.id = team_to_character.character_id) WHERE team_to_character.team_id = $1", team.ID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer charactersQuery.Close() // Закрываем запрос после использования

		var characters []Character
		for charactersQuery.Next() {
			var character Character
			err = charactersQuery.Scan(&character.ID, &character.Name, &character.Role, &character.Lore, &character.Img)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			characters = append(characters, character)
		}
		team.Characters = characters
		teams = append(teams, team)
	}
	fmt.Println(teams)
	return teams, nil
}
func getAllCharactersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	characters, err := getAllCharacters(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)

}

func getAllCharacters(db *sql.DB) ([]Character, error) {
	rows, err := db.Query("SELECT * FROM characters")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var characters []Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&character.ID, &character.Name, &character.Role, &character.Lore, &character.Img)
		characters = append(characters, character)
		if err != nil {
			return nil, err
		}
	}
	return characters, nil

}

func handleRequest() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/Pages/", http.StripPrefix("/Pages/", http.FileServer(http.Dir("./Pages/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/api/characters", getAllCharactersHandler)
	http.HandleFunc("/api/teams", getAllTeamsHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func main() {
	db, err := connectDB()
	if err != nil {
		fmt.Println("Не удалось подключиться к базе данных. Завершение работы...")
		return
	}
	defer db.Close() // Закрываем соединение только в main

	handleRequest()
}
