package utils

import (
	"crypto/tls"
	"errors"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

type LDAP_CONFIG struct {
	Addr       string   `json:"addr"`
	BaseDn     string   `json:"baseDn"`
	BindDn     string   `json:"bindDn`
	BindPass   string   `json:"bindPass"`
	AuthFilter string   `json:"authFilter"`
	Attributes []string `json:"attributes"`
	TLS        bool     `json:"tls"`
	StartTLS   bool     `json:"startTLS"`
	Conn       *ldap.Conn
}

type LDAP_RESULT struct {
	DN         string              `json:"dn"`
	Attributes map[string][]string `json:"attributes"`
}

func (lc *LDAP_CONFIG) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

func (lc *LDAP_CONFIG) Connect() (err error) {
	if lc.TLS {
		lc.Conn, err = ldap.DialTLS("tcp", lc.Addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		lc.Conn, err = ldap.Dial("tcp", lc.Addr)
	}
	if err != nil {
		return err
	}
	if !lc.TLS && lc.StartTLS {
		err = lc.Conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			lc.Conn.Close()
			return err
		}
	}

	err = lc.Conn.Bind(lc.BindDn, lc.BindPass)
	if err != nil {
		lc.Conn.Close()
		return err
	}
	return err
}

func (lc *LDAP_CONFIG) Auth(username, password string) (err error) {
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.AuthFilter, username), // The filter to apply
		lc.Attributes,                        // A list attributes to retrieve
		nil,
	)
	sr, err := lc.Conn.Search(searchRequest)
	if err != nil {
		return
	}
	if len(sr.Entries) == 0 {
		err = errors.New("Cannot find such user")
		return
	}
	if len(sr.Entries) > 1 {
		err = errors.New("Multi users in search")
		return
	}
	err = lc.Conn.Bind(sr.Entries[0].DN, password)
	if err != nil {
		return
	}
	//Rebind as the search user for any further queries
	err = lc.Conn.Bind(lc.BindDn, lc.BindPass)
	if err != nil {
		return
	}
	return
}

func LDAP_Auth(lc *LDAP_CONFIG, username, password string) (err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		return
	}
	err = lc.Auth(username, password)
	return

}
