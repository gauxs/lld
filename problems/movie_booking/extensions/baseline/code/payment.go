package code

type Pricer interface {
	Price(booking *Booking) *Price
}

type Price struct {
	base float64
}

type StandardPricing struct {
}

func (sp *StandardPricing) Price(booking *Booking) *Price {
	return nil
}
