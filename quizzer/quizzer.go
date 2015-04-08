// perform as a quizzer to root server
package quizzer

import "github.com/chzyer/adrs/dns"

type Quizzer struct {
	InChan chan *dns.DNSMessage
}

func NewQuizzer() {

}
