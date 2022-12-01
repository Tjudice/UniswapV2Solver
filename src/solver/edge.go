package solver

type Edge struct {
}

type Converter interface {
	// Need a function to compute the output amount given an input amount.
	// This must account for gas fees.
	// Will probably also have to account for direction? Not sure
	Cost()

	// Need a function to show the path of the edge
	// this is the relative input amounts, output amounts, and gas fees.
	// for now purely display, but can later switch to utilize call data for chaining requests together
	GetPathString()
}
