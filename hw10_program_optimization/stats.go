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
const DefaultDomainStatSize = 700

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	i := 0
	userEmail := UserEmail{}
	domainStat := make(DomainStat, DefaultDomainStatSize)
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

	if !strings.Contains(m.Email, "@") {
		return false
	}

	fullDomain := strings.ToLower(strings.SplitN(m.Email, "@", 2)[1])
	m.setEmailDomain(fullDomain)

	if !strings.Contains(fullDomain, ".") {
		return false
	}

	return strings.SplitN(fullDomain, ".", 2)[1] == domain
}

func (m *UserEmail) setEmailDomain(domain string) {
	m.emailDomain = domain
}

func (m *UserEmail) GetEmailDomain() string {
	return m.emailDomain
}
