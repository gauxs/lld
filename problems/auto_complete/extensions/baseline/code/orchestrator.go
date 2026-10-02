package code

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gauxs/lld/problems/auto_complete/extensions/baseline/code/enum"
)

func Execute() {
	ac := NewAutocomplete()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Autocomplete System CLI ===")
	fmt.Println("Commands:")
	fmt.Println("  add <word>          - Add a new word to the system")
	fmt.Println("  search <prefix>     - Search similar words matching the prefix")
	fmt.Println("  mode <order_type>   - Set search order (e.g., frequency)")
	fmt.Println("  exit                - Exit the program")
	fmt.Println("===============================")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Split command action from arguments
		parts := strings.Fields(line)
		command := strings.ToLower(parts[0])
		args := parts[1:]

		switch command {
		case "exit", "quit":
			fmt.Println("Exiting Autocomplete CLI. Goodbye!")
			return

		case "add":
			if len(args) < 1 {
				fmt.Println("❌ Error: Missing word. Usage: add <word>")
				continue
			}
			word := args[0]
			if err := ac.AddWord(word); err != nil {
				fmt.Printf("❌ Error adding word: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully added word: %q\n", word)
			}

		case "search":
			// If no argument is provided, default to an empty string to return top N overall
			prefix := ""
			if len(args) >= 1 {
				prefix = args[0]
			}

			results, err := ac.SearchSimilar(prefix)
			if err != nil {
				fmt.Printf("❌ Search failed: %v\n", err)
				continue
			}

			if len(results) == 0 {
				fmt.Println("🔍 No matching similarities found.")
			} else {
				fmt.Printf("🔍 Top Results: %s\n", strings.Join(results, ", "))
			}

		case "mode":
			if len(args) < 1 {
				fmt.Println("❌ Error: Specify a mode. Usage: mode frequency")
				continue
			}

			var targetOrder enum.SearchOrder
			modeStr := strings.ToLower(args[0])

			switch modeStr {
			case "frequency":
				targetOrder = enum.SEARCHORDER_FREQUENCY
			default:
				fmt.Printf("❌ Unknown mode %q. Defaulting or failing.\n", modeStr)
				continue
			}

			if err := ac.SetSearchOrder(targetOrder); err != nil {
				fmt.Printf("❌ Failed to change search order: %v\n", err)
			} else {
				fmt.Printf("⚙️  Search mode updated to: %s\n", modeStr)
			}

		default:
			fmt.Printf("❌ Unknown command: %q. Try 'add', 'search', 'mode', or 'exit'.\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
	}
}
