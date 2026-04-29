package styleguide

import (
	"fmt"
	"strings"
)

/* ---------------------------------------- Parser State ---------------------------------------- */

type blockMetadataParser struct {
	fields metadataFields
	field  metadataField
	value  string
	seen   map[metadataField]bool
}

/* ------------------------------------------- Parsing ------------------------------------------ */

func parseBlockMetadata(payload string) (fields metadataFields, err error) {
	parser := blockMetadataParser{seen: make(map[metadataField]bool)}

	for rawLine := range strings.SplitSeq(payload, "\n") {
		if err := parser.parseLine(rawLine); err != nil {
			return metadataFields{}, err
		}
	}

	if err := parser.finishField(); err != nil {
		return metadataFields{}, err
	}

	return parser.fields, nil
}

func (p *blockMetadataParser) parseLine(raw string) (err error) {
	raw = strings.TrimRight(raw, " \t\r")
	line := strings.TrimSpace(raw)
	if line == "" {
		return nil
	}

	indented := strings.TrimLeft(raw, " \t") != raw
	if indented {
		return p.appendContinuation(line)
	}

	name, value, hasAssignment := strings.Cut(line, "=")
	if !hasAssignment {
		if p.field == metadataFieldReason {
			return p.appendContinuation(line)
		}

		return fmt.Errorf("malformed style metadata comment near %q", line)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("malformed style metadata comment near %q", line)
	}

	field, found := parseMetadataField(name)
	if !found {
		return unknownMetadataField(metadataField(name))
	}

	return p.startField(field, strings.TrimSpace(value))
}

/* ----------------------------------------- Field State ---------------------------------------- */

func (p *blockMetadataParser) appendContinuation(line string) (err error) {
	if p.field != metadataFieldReason {
		return fmt.Errorf("malformed style metadata comment near %q", line)
	}

	if p.value != "" {
		p.value += " "
	}

	p.value += line
	return nil
}

func (p *blockMetadataParser) startField(field metadataField, value string) (err error) {
	if err := p.finishField(); err != nil {
		return err
	}

	p.field = field
	p.value = value
	return nil
}

func (p *blockMetadataParser) finishField() (err error) {
	if p.field == "" {
		return nil
	}

	if p.seen[p.field] {
		return fmt.Errorf("duplicate %q in style metadata comment", p.field)
	}

	if err := p.fields.setField(p.field, p.value); err != nil {
		return err
	}

	p.seen[p.field] = true
	p.field, p.value = "", ""
	return nil
}
