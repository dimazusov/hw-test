package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type UserEmail struct {
	Email       string
	emailDomain string
}

type DomainStat map[string]int

const ErrMessageReadLine = "cannot read line"
const ErrMessageUnmarshal = "cannot unmarshal"

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	i := 0
	userEmail := UserEmail{}
	domainStat := make(DomainStat, 700)
	bufReader := bufio.NewReader(r)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return domainStat, nil
			}
			return nil, errors.Wrap(err, ErrMessageReadLine)
		}

		if err = userEmail.UnmarshalJSON(line); err != nil {
			return nil, errors.Wrap(err, ErrMessageUnmarshal)
		}

		if userEmail.IsEmailHasDomain(domain) {
			domainStat[userEmail.GetEmailDomain()]++
		}

		i++
	}
}

func (m *UserEmail) IsEmailHasDomain(domain string) bool {
	if len(m.Email) == 0 {
		return false
	}

	fullDomain := strings.ToLower(strings.SplitN(m.Email, "@", 2)[1])
	m.setEmailDomain(fullDomain)

	return strings.SplitN(fullDomain, ".", 2)[1] == domain
}

func (m *UserEmail) setEmailDomain(domain string) {
	m.emailDomain = domain
}

func (m *UserEmail) GetEmailDomain() string {
	return m.emailDomain
}
