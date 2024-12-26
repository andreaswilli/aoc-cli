package trigger

type Trigger interface {
  Listen() chan bool
}
