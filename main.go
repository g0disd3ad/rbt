package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/g0disd3ad/rbt/internal/api"
	"github.com/g0disd3ad/rbt/internal/dictionary"
	"github.com/g0disd3ad/rbt/internal/storage"
)

func isStrictEnglish(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if r == '-' || r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			continue
		}
		return false
	}
	return true
}

func isStrictRussian(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if r == '-' || r == ' ' || (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'ё' || r == 'Ё' {
			continue
		}
		return false
	}
	return true
}

func printMenu() {
	fmt.Println("\n=== DICTIONARY MENU ===")
	fmt.Println("1. Add translation (eng - rus)")
	fmt.Println("2. Find the word (eng)")
	fmt.Println("3. Remove the word (eng)")
	fmt.Println("4. Print the dictionary")
	fmt.Println("5. Save to txt file")
	fmt.Println("6. Save to PostgreSQL database")
	fmt.Println("7. Checking the tree (height and validity)")
	fmt.Println("0. Exit")
	fmt.Print("Option: ")
}

func main() {
	rbtStorage := dictionary.NewRBTStorage()
	dict := dictionary.NewDictionary(rbtStorage)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("~--- Initialising the dictionary ---~")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var pg *storage.PostgresStorage
	var dbErr error

	if dbPassword == "" {
		fmt.Println("Warning: DB_PASSWORD environment variable is missing. Database storage is disabled.")
	} else {
		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		if dbUser == "" {
			dbUser = "postgres"
		}
		if dbName == "" {
			dbName = "dict_db"
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
		pg, dbErr = storage.NewPostgresStorage(dsn)
		if dbErr != nil {
			fmt.Printf("Warning: Could not connect to DB: %v\n", dbErr)
		}
	}

	if dbErr != nil {
		fmt.Printf("Warning: Could not connect to DB: %v\n", dbErr)
	} else {
		defer pg.Close()
		fmt.Print("Load data from PostgreSQL database? (y/n): ")
		if scanner.Scan() {
			ans := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if ans == "y" || ans == "yes" || ans == "lf" {
				err := pg.LoadToTree(dict.Insert)
				if err != nil {
					fmt.Printf("ERROR loading from DB: %v\n", err)
				}
			}
		}
	}

	fmt.Print("Enter the name of the file for loading or press Enter to skip: ")
	if scanner.Scan() {
		filename := strings.TrimSpace(scanner.Text())
		if filename != "" {
			err := dict.LoadFromFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			}
		}
	}
	apiServer := api.NewAPI(dict)
	apiServer.StartServer(":8080")
	fmt.Println("API Server is running on http://localhost:8080")
	for {
		printMenu()
		if !scanner.Scan() {
			break
		}
		inputStr := strings.TrimSpace(scanner.Text())
		if inputStr == "" {
			continue
		}

		switch inputStr {
		case "0":
			fmt.Println("Exiting the program...")
			return

		case "1":
			fmt.Print("Enter 'eng - rus': ")
			if scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, " - ", 2)

				if len(parts) == 2 {
					eng := strings.TrimSpace(parts[0])
					rus := strings.TrimSpace(parts[1])

					if !isStrictEnglish(eng) || !isStrictRussian(rus) {
						fmt.Println("  ERROR: Check the input language (without whitespaces).")
					} else {
						err := dict.Insert(eng, rus)
						if err != nil {
							fmt.Printf("  ERROR: %v\n", err)
						} else {
							fmt.Printf("  Word %s - %s is added.\n", eng, rus)
						}
					}
				} else {
					fmt.Println("  ERROR: Format is 'english word - russian word'.")
				}
			}

		case "2":
			fmt.Print("Enter english word for searching: ")
			if scanner.Scan() {
				eng := strings.TrimSpace(scanner.Text())
				translations, err := dict.Search(eng)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Translations: %s\n", strings.Join(translations, ", "))
				}
			}

		case "3":
			fmt.Print("Enter the english word for removing: ")
			if scanner.Scan() {
				eng := strings.TrimSpace(scanner.Text())
				err := dict.Remove(eng)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Word %s is removed.\n", eng)
				}
			}

		case "4":
			dict.Print()

		case "5":
			fmt.Print("Enter the name of the file for saving: ")
			if scanner.Scan() {
				saveName := strings.TrimSpace(scanner.Text())
				err := dict.SaveToFile(saveName)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
				}
			}

		case "6":
			if pg != nil {
				err := pg.SaveFromTree(dict)
				if err != nil {
					fmt.Printf("ERROR saving to DB: %v\n", err)
				}
			} else {
				fmt.Println("ERROR: No database connection established.")
			}

		case "7":
			height := dict.GetHeight()
			isValid := dict.IsValidTree()

			fmt.Println("\n~--- Checking the Red Black Tree ---~")
			fmt.Printf("Current maximum height of the tree: %d\n", height)
			if isValid {
				fmt.Println("Status of validation: OK")
			} else {
				fmt.Println("Status of validation: ERROR (Tree is not following the rules of Red Black Tree!)")
			}
			fmt.Println("~-----------------------------------------~")

		default:
			fmt.Println("Invalid option.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error of input: %v\n", err)
	}

	fmt.Println("The program has finished working.")
}
