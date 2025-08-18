package cmd

import "github.com/charmbracelet/lipgloss"

// Color palette for ctx CLI styling
var (
	// Header styles for main sections
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B35")).
			Bold(true)

	// Subheader styles for workflow names
	subHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0088CC")).
			Bold(true)

	// Command/example styles
	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0088CC"))

	// Description styles
	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))

	// Important text styles
	importantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AA00")).
			Bold(true)

	// Separator style
	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))

	// Environment variable name style
	envVarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)
)

// StyledHeader returns a styled header string
func StyledHeader(text string) string {
	return headerStyle.Render(text)
}

// StyledSubHeader returns a styled subheader string
func StyledSubHeader(text string) string {
	return subHeaderStyle.Render(text)
}

// StyledCommand returns a styled command string
func StyledCommand(text string) string {
	return commandStyle.Render(text)
}

// StyledDescription returns a styled description string
func StyledDescription(text string) string {
	return descStyle.Render(text)
}

// StyledImportant returns a styled important text string
func StyledImportant(text string) string {
	return importantStyle.Render(text)
}

// StyledSeparator returns a styled separator string
func StyledSeparator(text string) string {
	return separatorStyle.Render(text)
}

// StyledEnvVar returns a styled environment variable name
func StyledEnvVar(text string) string {
	return envVarStyle.Render(text)
}
