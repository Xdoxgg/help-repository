package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
	"strconv"
)

type IdToId struct {
	ID  int `json:"id"`
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
	Role     string `json:"role"`
}

type DataTeam struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Orientation string `json:"orientation"`
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

// /////////////////////////
func getAllTeamsDataHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	characters, err := getAllTeamsData(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)

}

func getAllTeamsData(db *sql.DB) ([]DataTeam, error) {
	rows, err := db.Query("SELECT * FROM teams")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var characters []DataTeam
	for rows.Next() {
		var character DataTeam
		err := rows.Scan(&character.ID, &character.Name, &character.Orientation)
		characters = append(characters, character)
		if err != nil {
			return nil, err
		}
	}
	return characters, nil

}

func getAllTeamsToDataHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	characters, err := getAllTeamsToData(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)

}

func getAllTeamsToData(db *sql.DB) ([]IdToId, error) {
	rows, err := db.Query("SELECT * FROM team_to_character")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var characters []IdToId
	for rows.Next() {
		var character IdToId
		err := rows.Scan(&character.ID, &character.Id1, &character.Id2)
		characters = append(characters, character)
		if err != nil {
			return nil, err
		}
	}
	return characters, nil

}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	characters, err := getAllUsers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)

}

func getAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var characters []User
	for rows.Next() {
		var character User
		err := rows.Scan(&character.ID, &character.Name, &character.Password, &character.Role)
		characters = append(characters, character)
		if err != nil {
			return nil, err
		}
	}
	return characters, nil

}

func deleteTeamToCharacterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполняем SQL-запрос на удаление
	_, err = db.Exec("DELETE FROM team_to_character WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func deleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполняем SQL-запрос на удаление
	_, err = db.Exec("DELETE FROM teams WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func deleteCharactersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выполняем SQL-запрос на удаление
	_, err = db.Exec(`
DELETE FROM characters WHERE id = $1
`, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		fmt.Println("err")
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("start")
	// Выполняем SQL-запрос на удаление
	fmt.Println(id)
	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func addCharacterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	img := r.URL.Query().Get("img")
	name := r.URL.Query().Get("name")
	role := r.URL.Query().Get("role")
	lore := r.URL.Query().Get("lore")
	// Выполняем SQL-запрос на удаление
	_, err = db.Exec(`
INSERT INTO characters (name, role, lore, talents_build_emblems)
VALUES ($1, $2, $3, $4);

`, name, role, lore, img)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//teams /api/teams_data/update?name=newteam&orientation=s&id=
//teasm_items/api/team_to_character/update?team_id=1&character_id=1&id=
//users/api/users/update?name=as&password=s&id=

func addTeamDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	name := r.URL.Query().Get("name")
	orientation := r.URL.Query().Get("orientation")

	_, err = db.Exec(`INSERT INTO teams (name, orientation)
VALUES ($1, $2);`, name, orientation)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addTeamLoCharacterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("err")
		return
	}
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	tId := r.URL.Query().Get("team_id")
	cId := r.URL.Query().Get("character_id")

	_, err = db.Exec(`INSERT INTO team_to_character (team_id, character_id)
VALUES ($1, $2);`, tId, cId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	db, err := connectDB()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")

	_, err = db.Exec(`
INSERT INTO users (username, password, role) VALUES
    ($1, $2, $3);`, name, password, 0)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func handleRequest() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/Pages/", http.StripPrefix("/Pages/", http.FileServer(http.Dir("./Pages/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/api/characters", getAllCharactersHandler)
	http.HandleFunc("/api/teams", getAllTeamsHandler)
	////////////////

	http.HandleFunc("/api/teams_data", getAllTeamsDataHandler)
	http.HandleFunc("/api/team_to_character", getAllTeamsToDataHandler)
	http.HandleFunc("/api/users", getAllUsersHandler)
	http.HandleFunc("/api/team_to_character_delete", deleteTeamToCharacterHandler)
	http.HandleFunc("/api/teams_data_delete", deleteTeamHandler)
	http.HandleFunc("/api/characters_delete", deleteCharactersHandler)
	http.HandleFunc("/api/users_delete", deleteUsersHandler)
	http.HandleFunc("/api/characters_add", addCharacterHandler)
	http.HandleFunc("/api/users_add", addUsersHandler)
	http.HandleFunc("/api/team_to_character_add", addTeamLoCharacterHandler)
	http.HandleFunc("/api/teams_data_add", addTeamDataHandler)

	////////////
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
