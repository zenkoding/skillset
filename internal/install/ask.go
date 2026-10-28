package install

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Ask asks which side wins for each unresolved settings conflict, reading
// answers from in and writing prompts to out.
func (p *Plan) Ask(in io.Reader, out io.Writer) error {
	r := bufio.NewReader(in)
	for i := range p.Copies {
		c := &p.Copies[i]
		for _, k := range c.Unresolved() {
			for {
				fmt.Fprintf(out, "%s %s\n  repo:  %s\n  local: %s\nkeep [r]epo or [l]ocal? ", c.Path, k.Key, k.Repo, k.Local)
				line, err := r.ReadString('\n')
				answer := strings.ToLower(strings.TrimSpace(line))
				if answer == "r" || answer == "l" {
					if err := p.Resolve(c.Path, k.Key, answer == "l"); err != nil {
						return err
					}
					break
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
