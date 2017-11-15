package genetic
import(
"fmt"
)

func Red(text string) string{
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

func Blue(text string) string{
	return fmt.Sprintf("\033[34m%s\033[0m", text)
}

func White(text string) string{
	return fmt.Sprintf("\033[1;37m%s\033[0m", text)
}

func Yellow(text string) string{
	return fmt.Sprintf("\033[33m%s\033[0m", text)
}