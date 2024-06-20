# Comments are ignored 

Scenario: Home page accessed
    When the home page is accessed
    Then the page contains `Learning with Flashcards.`

Scenario: Test deck accessed
    When the test deck is accessed
	Then the deck has the question `What is the meaning of life?`
