package pgsql

import (
	models "effectivemobiletest/internal/models/db-models"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

func (s *Storage) CreateUser(user models.UserIn) error {
	qstr := fmt.Sprintf(`
		INSERT INTO users (FName, SName, PName, GenderId, Age, RegionId)
		VALUES ('%s', '%s', '%s', %d, %d, %d)
	`, user.FName, user.SName, user.PName, user.GenderId, user.Age, user.RegionId)

	_, err := s.db.Exec(qstr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ReadUsers(filters map[string]string, limit int, offset int) ([]models.UserOut, error) {
	qstr := s.generateGetQueryString(filters, limit, offset)

	rows, err := s.db.Query(qstr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users = make([]models.UserOut, 0)
	var user models.UserOut

	for rows.Next() {
		err = rows.Scan(&user.UserId, &user.FName, &user.SName, &user.PName, &user.GenderName, &user.Age, &user.RegionCode)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) UpdateUser(user models.UserIn) error {
	qstr := fmt.Sprintf(`UPDATE Users SET
		FName = '%s', 
		SName = '%s', 
		PName = '%s', 
		GenderId = %d, 
		Age = %d, 
		RegionId = %d
		WHERE UserId = %d`,
		user.FName, user.SName, user.PName, user.GenderId, user.Age, user.RegionId, user.UserId)

	_, err := s.db.Exec(qstr)
	if err != nil {
		return err
	}

	//Лог на изменение данных
	return nil
}

func (s *Storage) DeleteUser(userId int) error {
	qstr := fmt.Sprintf(`DELETE FROM Users WHERE UserId=%d`, userId)
	_, err := s.db.Exec(qstr)
	if err != nil {
		return nil
	}
	// Добавить лог на удаление данных
	return nil
}

func (s *Storage) ReadRegion(regionCode string) (*models.Region, error) {
	qstr := fmt.Sprintf("SELECT * FROM Regions WHERE regionCode = '%s'", regionCode)
	row := s.db.QueryRow(qstr)
	var region models.Region
	err := row.Scan(&region.RegionId, &region.RegionCode)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (s *Storage) WriteRegion(regionCode string) error {
	qstr := fmt.Sprintf("INSERT INTO Regions VALUES (%s)", regionCode)
	_, err := s.db.Exec(qstr)
	if err != nil {
		return err
	}
	// Лог о записи данных
	return nil
}

func (s *Storage) GetGenders() ([]models.Gender, error) {
	qstr := "SELECT * FROM Genders"
	rows, err := s.db.Query(qstr)
	if err != nil {
		return nil, err
	}

	genders := make([]models.Gender, 0)
	var gender models.Gender
	for rows.Next() {
		err = rows.Scan(&gender.GenderId, &gender.GenderName)
		if err != nil {
			slog.Warn("Row doesn't scanning", slog.String("error", err.Error()))
			continue
		}
		genders = append(genders, gender)
	}
	return genders, nil
}

func (s *Storage) GetGenderId(gender string) (int, error) {
	var genderId int
	checkQuery := fmt.Sprintf("SELECT GenderId FROM Genders WHERE GenderName = '%s'", gender)
	req := s.db.QueryRow(checkQuery)
	err := req.Scan(&genderId)

	if err != nil {
		return 0, err
	}
	return genderId, nil
}

func (s *Storage) RegionCodeChekingAndCreating(code string) (int, error) {
	var regionId int
	code = strings.ToUpper(code)

	checkQuery := fmt.Sprintf("SELECT RegionId FROM Regions WHERE RegionCode = '%s'", code)
	req := s.db.QueryRow(checkQuery)
	err := req.Scan(&regionId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q := fmt.Sprintf(`INSERT INTO Regions (RegionCode) VALUES ('%s')`, code)
			_, err = s.db.Exec(q)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	req = s.db.QueryRow(checkQuery)
	err = req.Scan(&regionId)

	if err != nil {
		return 0, err
	}

	return regionId, nil
}

func (s *Storage) GetRandomRegionId() (int, error) {
	qstr := "SELECT RegionId FROM Regions ORDER BY RANDOM() LIMIT 1"
	row := s.db.QueryRow(qstr)
	var reg int
	err := row.Scan(&reg)

	if err != nil {
		return 0, err
	}
	return reg, nil

}

func (s *Storage) generateGetQueryString(filters map[string]string, limit int, offset int) string {
	qstr :=
		`SELECT UserId, FName, SName, PName, GenderName, Age,RegionCode
FROM users 
INNER JOIN regions ON users.RegionId = regions.RegionId
INNER JOIN genders ON users.GenderId = genders.GenderId
`
	if len(filters) > 0 {
		qstr += "WHERE "
		remParam := len(filters)
		for k, v := range filters {
			if r := regexp.MustCompile(`^\d+$`); r.Match([]byte(v)) {
				qstr += fmt.Sprintf("%s = %s", k, v)
			} else {
				qstr += fmt.Sprintf("%s = '%s'", k, v)
			}
			remParam--
			if remParam > 0 {
				qstr += " AND "
			}
		}
	}
	if limit > 0 {
		qstr += fmt.Sprintf(" AND UserId > %d", offset)
		qstr += fmt.Sprintf("\nORDER BY UserId ASC \nLIMIT %d", limit)
	}
	return qstr
}
