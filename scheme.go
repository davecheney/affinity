/*
   Affinity - Private groups as a service
   Copyright (C) 2014  Canonical, Ltd.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, version 3.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package affinity

import (
	"net/url"

	"code.google.com/p/gopass"
)

type PasswordProvider interface {
	Password() (string, error)
}

type PasswordPrompter struct{}

func (pp *PasswordPrompter) Password() (string, error) {
	return gopass.GetPass("Password: ")
}

type SchemeAuthorizer interface {
	// Auth creates the authorization parameters for the given identity. Other parameters (passphrases, private keys, etc.) may be used as factors in creating them for various schemes.
	Auth(id string) (values url.Values, err error)
}

type SchemeValidator interface {
	// Validate checks the authorization parameters are valid. If so, returns the
	// qualified user ID which created it.
	Validate(values url.Values) (id string, err error)
}

type Scheme interface {
	Authorizer() SchemeAuthorizer
	Name() string
	Validator() SchemeValidator
}

type SchemeMap map[string]Scheme

func (sm SchemeMap) Register(scheme Scheme) {
	sm[scheme.Name()] = scheme
}
