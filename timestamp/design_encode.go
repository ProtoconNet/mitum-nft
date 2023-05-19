package timestamp

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	svc string,
	prjs []string,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)
	de.service = extensioncurrency.ContractID(svc)
	de.projects = prjs

	return nil
}
