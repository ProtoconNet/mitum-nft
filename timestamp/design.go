package timestamp

import (
	"bytes"
	"sort"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var DesignHint = hint.MustNewHint("mitum-timestamp-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	service  extensioncurrency.ContractID
	projects []string
}

func NewDesign(service extensioncurrency.ContractID, projects ...string) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		service:    service,
		projects:   projects,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.service,
	); err != nil {
		return err
	}

	return nil
}

func (de Design) Bytes() []byte {
	length := 1
	bytesArray := make([][]byte, length+len(de.projects))
	bytesArray[0] = de.service.Bytes()

	sort.Slice(de.projects, func(i, j int) bool {
		return bytes.Compare([]byte(de.projects[j]), []byte(de.projects[i])) < 0
	})

	for i := range de.projects {
		bytesArray[i+length] = []byte(de.projects[i])
	}

	return util.ConcatBytesSlice(bytesArray...)
}

func (de Design) Hash() util.Hash {
	return de.GenerateHash()
}

func (de Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(de.Bytes())
}

func (de Design) Service() extensioncurrency.ContractID {
	return de.service
}

func (de Design) Projects() []string {
	return de.projects
}

func (de *Design) AddProject(project string) {
	for i := range de.projects {
		if de.projects[i] == project {
			return
		}
	}
	projects := append(de.projects, project)
	de.projects = projects
}

func (de Design) Equal(cd Design) bool {
	if de.service != cd.service {
		return false
	}

	if len(de.projects) != len(cd.projects) {
		return false
	}

	sort.Slice(de.projects, func(i, j int) bool {
		return bytes.Compare([]byte(de.projects[i]), []byte(de.projects[j])) < 0
	})

	for i := range de.projects {
		if de.projects[i] != cd.projects[i] {
			return false
		}
	}

	return true
}
