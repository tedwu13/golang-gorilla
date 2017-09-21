package funding


type Fund struct {
    // balance is unexported since it is lowercase
    // if it is uppercase it can be exported
    balance int 
}

// A regular function returning a pointer to a fund
func NewFund(initialBalance int) *Fund {
    // Return a pointer to a new struct without worrying about whehter is it on stack or heap. Go figures that out for us
    return &Fund {
        balance : initialBalance,
    }
}

// Methods starts with a *receiver* in this cas it is a Fund pointer
func (fund *Fund) Balance() int {
    return fund.balance
}

func (fund *Fund) Withdraw(amount int) {
    fund.balance -= amount
}


// 2 billion iterations under 3 seconds
