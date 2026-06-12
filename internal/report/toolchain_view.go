package report

type ToolchainView struct {
	Result   ToolchainResult
	AllValid bool
}

func NewToolchainView(result ToolchainResult) (view ToolchainView) {
	view = ToolchainView{
		Result:   result,
		AllValid: true,
	}
	for _, status := range result.Statuses {
		if status.Valid {
			continue
		}

		view.AllValid = false
		break
	}

	return view
}
