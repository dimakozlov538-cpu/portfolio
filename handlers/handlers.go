package handlers

import (
	"html/template"
	"log"
	"net/http"
	"portfolio/db"
	"portfolio/models"

	"github.com/lib/pq"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		log.Println("КРИТИЧЕСКАЯ ОШИБКА: db.DB is nil!")
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	data := models.PageData{}

	// 1. Профиль
	db.DB.QueryRow("SELECT id, name, title, photo_url FROM profile LIMIT 1").
		Scan(&data.Profile.ID, &data.Profile.Name, &data.Profile.Title, &data.Profile.PhotoURL)

	// 2. Обо мне
	db.DB.QueryRow("SELECT id, title, description FROM about LIMIT 1").
		Scan(&data.About.ID, &data.About.Title, &data.About.Description)

	// 3. Навыки
	rows, _ := db.DB.Query("SELECT id, name, category FROM skills ORDER BY category")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var s models.Skill
			rows.Scan(&s.ID, &s.Name, &s.Category)
			data.Skills = append(data.Skills, s)
		}
	}

	// 4. ПРОЕКТЫ (исправленная, надёжная версия)
	query := `
		SELECT id, title, description, 
		       COALESCE(image_url, ''), 
		       COALESCE(repo_url, ''), 
		       COALESCE(demo_url, ''), 
		       COALESCE(tech_stack, '{}') 
		FROM projects`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Ошибка запроса projects: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var p models.Project
			var stack pq.StringArray

			err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.ImageURL, &p.RepoURL, &p.DemoURL, &stack)
			if err != nil {
				log.Printf("Ошибка сканирования проекта: %v", err)
				continue
			}

			p.TechStack = stack

			log.Printf("✅ DEBUG: Загружен проект '%s' | RepoURL: '%s' | Tech: %v", p.Title, p.RepoURL, p.TechStack)

			data.Projects = append(data.Projects, p)
		}
	}

	// 5. Хакатоны
	rows, _ = db.DB.Query("SELECT id, title, description, date, result FROM hackathons ORDER BY date DESC")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var h models.Hackathon
			rows.Scan(&h.ID, &h.Title, &h.Description, &h.Date, &h.Result)
			data.Hackathons = append(data.Hackathons, h)
		}
	}

	// 6. Награды
	rows, _ = db.DB.Query("SELECT id, title, description, date FROM awards ORDER BY date DESC")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var a models.Award
			rows.Scan(&a.ID, &a.Title, &a.Description, &a.Date)
			data.Awards = append(data.Awards, a)
		}
	}

	// 7. Образование
	rows, _ = db.DB.Query("SELECT id, degree, institution, description, start_date, end_date FROM education")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var e models.Education
			rows.Scan(&e.ID, &e.Degree, &e.Institution, &e.Description, &e.StartDate, &e.EndDate)
			data.Education = append(data.Education, e)
		}
	}

	// 8. Контакты
	rows, _ = db.DB.Query("SELECT id, type, value, url FROM contacts")
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var c models.Contact
			rows.Scan(&c.ID, &c.Type, &c.Value, &c.URL)
			data.Contacts = append(data.Contacts, c)
		}
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Ошибка парсинга шаблона: %v", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона: %v", err)
	}
}
