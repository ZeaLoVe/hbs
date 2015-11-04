package db

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"log"
	"time"
)

func QueryHosts() (map[string]int, map[string]string, error) {
	m := make(map[string]int)
	m2 := make(map[string]string)

	sql := "select id, hostname ,ip from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, m2, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id       int
			hostname string
			ip       string
		)

		err = rows.Scan(&id, &hostname, &ip)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m[hostname] = id
		m2[hostname] = ip
	}

	return m, m2, nil
}

func QueryMonitoredHosts() (map[int]*model.Host, error) {
	hosts := make(map[int]*model.Host)
	now := time.Now().Unix()
	sql := fmt.Sprintf("select id, hostname from host where maintain_begin > %d or maintain_end < %d", now, now)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return hosts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := model.Host{}
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			log.Println("WARN:", err)
			continue
		}
		hosts[t.Id] = &t
	}

	return hosts, nil
}

//Get endpoint by ip
func QueryEndpoint(ip string) (string, error) {
	sql := fmt.Sprintf("select hostname from host where ip = '%v' order by update_at DESC", ip)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return "", err
	}

	defer rows.Close()
	for rows.Next() {
		var hostname string

		err = rows.Scan(&hostname)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		return hostname, nil
	}
	return "", fmt.Errorf("Can't get endpoint")
}
