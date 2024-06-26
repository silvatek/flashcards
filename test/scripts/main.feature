# Post-deployment tests for Flashcards webapp

Feature: Home page 

    Scenario: Home page accessed
        When the home page is accessed
        Then the page contains `Learning with Flashcards.`
        And  the page contains `You can open an existing deck`

Feature: Deck page

    Scenario: Test deck accessed
        When the test deck is accessed
        Then the page contains `Test flashcard deck`
        And  the deck has the question `What is the meaning of life?`
