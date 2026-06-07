package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (ss *SchoolSystem) LoadAllData() error {
	if ss.Users == nil { ss.Users = make(map[string]User) }
	if ss.Results == nil { ss.Results = make(map[string]StudentResult) }
	if ss.Notices == nil { ss.Notices = make(map[string]Notice) }
	if ss.Timetables == nil { ss.Timetables = make(map[string]TimetableItem) }
	if ss.PendingOTPs == nil { ss.PendingOTPs = make(map[string]OTP) }

	// 1. LOAD USERS
	uFile, err := os.Open(ss.UserFile)
	if err == nil {
		defer uFile.Close()
		scanner := bufio.NewScanner(uFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" { continue }
			p := strings.Split(line, ":")
			if len(p) == 7 {
				if _, exists := ss.Users[p[0]]; !exists {
					ss.Users[p[0]] = User{
						ID: p[0], Surname: p[1], MiddleName: p[2], FirstName: p[3], Password: p[4], Email: p[5], Role: p[6],
					}
				}
			}
		}
		if scanErr := scanner.Err(); scanErr != nil { return scanErr }
	}

	// 2. LOAD RESULTS
	rFile, err := os.Open(ss.ResultFile)
	if err == nil {
		defer rFile.Close()
		scanner := bufio.NewScanner(rFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" { continue }
			p := strings.Split(line, ":")
			if len(p) == 4 {
				key := p[0] + "_" + p[1]
				if _, exists := ss.Results[key]; !exists {
					units, _ := strconv.Atoi(p[3])
					ss.Results[key] = StudentResult{
						StudentID: p[0], CourseCode: p[1], Grade: p[2], CreditUnits: units,
					}
				}
			}
		}
		if scanErr := scanner.Err(); scanErr != nil { return scanErr }
	}
	return nil
}

func (ss *SchoolSystem) RewriteUsers() {
	file, err := os.Create(ss.UserFile)
	if err != nil { return }
	defer file.Close()
	for _, u := range ss.Users {
		fmt.Fprintf(file, "%s:%s:%s:%s:%s:%s:%s\n", u.ID, u.Surname, u.MiddleName, u.FirstName, u.Password, u.Email, u.Role)
	}
}

func (ss *SchoolSystem) RewriteResults() {
	file, _ := os.Create(ss.ResultFile); defer file.Close()
	for _, r := range ss.Results { fmt.Fprintf(file, "%s:%s:%s:%d\n", r.StudentID, r.CourseCode, r.Grade, r.CreditUnits) }
}

func (ss *SchoolSystem) RewriteNotices() {
	file, _ := os.Create(ss.NoticeFile); defer file.Close()
	for _, n := range ss.Notices { fmt.Fprintf(file, "%s|%s|%s|%s\n", n.ID, n.Author, n.Content, n.Timestamp.Format(time.RFC3339)) }
}

func (ss *SchoolSystem) RewriteTimetables() {
	file, _ := os.Create(ss.TimetableFile); defer file.Close()
	for _, t := range ss.Timetables { fmt.Fprintf(file, "%s:%s:%s:%s:%s\n", t.ID, t.Type, t.CourseCode, t.DateTime, t.Venue) }
}