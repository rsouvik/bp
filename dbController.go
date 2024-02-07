package main

type DbCtrl struct {
}

func NewDbCtrl() *DbCtrl {
	return &DbCtrl{}
}

func (p *DbCtrl) Run(s *SharedExtConn) error {

	return nil

}
