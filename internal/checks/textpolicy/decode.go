package textpolicy

import "fmt"

func configSection(
	pack map[string]any,
	key string,
	field string,
) (section map[string]any, err error) {
	value, found := pack[key]
	if !found {
		return nil, nil
	}

	section, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a table", field)
	}

	return section, nil
}

func decodeSectionHeaderConfig(
	section map[string]any,
) (config SectionHeaderConfig, err error) {
	if err = rejectUnknownFields(
		section,
		"packs.text.section_headers",
		"large_min_lines",
		"short_max_lines",
		"max_header_count",
		"generic_names",
		"structural_names",
	); err != nil {
		return SectionHeaderConfig{}, err
	}

	config.LargeMinLines, err = intField(
		section,
		"large_min_lines",
		"packs.text.section_headers.large_min_lines",
	)
	if err != nil {
		return SectionHeaderConfig{}, err
	}

	config.ShortMaxLines, err = intField(
		section,
		"short_max_lines",
		"packs.text.section_headers.short_max_lines",
	)
	if err != nil {
		return SectionHeaderConfig{}, err
	}

	config.MaxHeaderCount, err = intField(
		section,
		"max_header_count",
		"packs.text.section_headers.max_header_count",
	)
	if err != nil {
		return SectionHeaderConfig{}, err
	}

	config.GenericNames, err = stringList(
		section,
		"generic_names",
		"packs.text.section_headers.generic_names",
	)
	if err != nil {
		return SectionHeaderConfig{}, err
	}

	config.StructuralNames, err = stringList(
		section,
		"structural_names",
		"packs.text.section_headers.structural_names",
	)
	if err != nil {
		return SectionHeaderConfig{}, err
	}

	return config, nil
}
