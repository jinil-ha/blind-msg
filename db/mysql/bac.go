package mysql

import (
	_ "github.com/go-sql-driver/mysql" // using mysql
	"github.com/kataras/golog"

	svc "github.com/jinil-ha/blind-msg/service"
)

// GetBAC select bac info from bac table
func GetBAC(service string, userid string, bac *string) error {
	s := svc.GetCode(service)
	golog.Debugf("select bac: %d %s", s, userid)

	selDB, err := database.Query("SELECT bac FROM bac WHERE service = ? AND user_id = ?", s, userid)
	if err != nil {
		golog.Warnf("DB Error : %s", err)
		return err
	}

	if selDB.Next() {
		var v string
		err = selDB.Scan(&v)
		if err != nil {
			return err
		}
		*bac = v
	} else {
		*bac = ""
	}

	return nil
}

// SetBAC insert bac info to bac table
func SetBAC(service string, userid string, bac string) error {
	s := svc.GetCode(service)
	golog.Debugf("insert bac: %d %s %s", s, userid, bac)

	ins, err := database.Prepare("INSERT INTO bac (service, user_id, bac) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	ins.Exec(s, userid, bac)

	return nil
}

// CheckBAC ...
func CheckBAC(bac string) error {
	// TODO
	return nil
}

// GetUserInfo get user bac info from DB
func GetUserInfo(bac string, service *string, userID *string) error {
	// s := getServiceCode(service)
	golog.Debugf("select userid: bac(%s)", bac)

	sel, err := database.Query("SELECT service, user_id FROM bac WHERE bac = ?", bac)
	if err != nil {
		return err
	}

	if sel.Next() {
		var v1 int
		var v2 string
		err = sel.Scan(&v1, &v2)
		if err != nil {
			return err
		}
		*service = svc.GetName(v1)
		*userID = v2
	}

	return nil
}
