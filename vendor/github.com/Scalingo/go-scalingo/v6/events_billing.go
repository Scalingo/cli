package scalingo

import (
	"fmt"

	"github.com/Scalingo/go-scalingo/v6/billing"
)

type EventAddCreditType struct {
	Event
	TypeData EventAddCreditTypeData `json:"type_data"`
}

func (ev *EventAddCreditType) String() string {
	return fmt.Sprintf(
		"%f€ of credit added to your account (%s)", ev.TypeData.Amount, ev.TypeData.PaymentMethod,
	)
}

type EventAddCreditTypeData struct {
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
}

type EventAddPaymentMethodType struct {
	Event
	TypeData EventAddPaymentMethodTypeData `json:"type_data"`
}

func (ev *EventAddPaymentMethodType) String() string {
	if ev.TypeData.Profile.PaymentMethodType == billing.Stripe {
		p := ev.TypeData.Profile.Stripe
		return fmt.Sprintf("%s card ending with ...%s, expiring in %s", p.Brand, p.Last4, p.Exp)
	}

	return fmt.Sprintf("'%s' payment method added", ev.TypeData.Profile.PaymentMethodType)
}

type EventAddPaymentMethodTypeData struct {
	billing.Profile
}

type EventAddVoucherType struct {
	Event
	TypeData EventAddVoucherTypeData `json:"type_data"`
}

func (ev *EventAddVoucherType) String() string {
	return fmt.Sprintf("code: '%s'", ev.TypeData.Code)
}

type EventAddVoucherTypeData struct {
	Code string `json:"code"`
}
