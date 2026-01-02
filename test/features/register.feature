Feature: User Registration
  In order to access the system
  As a new user
  I want to register with my credentials

  Scenario: Register user with valid data
    Given I have valid user data:
      | name      | John Doe          |
      | email     | john@test.com     |
      | password  | password123       |
    When I send a POST request to "/users"
    Then I should receive a 201 status code
    And the response should contain a "token" field
    And the response should not contain an error message

  Scenario: Register user with existing email
    Given I have valid user data:
      | name      | John Doe          |
      | email     | existing@test.com |
      | password  | password123       |
    And the user "existing@test.com" is already registered
    When I send a POST request to "/users"
    Then I should receive a 400 status code
    And the response should contain error message "user already exists"

  Scenario: Register user with invalid email
    Given I have invalid user data:
      | name      | John Doe      |
      | email     | invalidemail  |
      | password  | password123   |
    When I send a POST request to "/users"
    Then I should receive a 400 status code
    And the response should contain error message "email is invalid"

  Scenario: Register user with short password
    Given I have invalid user data:
      | name      | John Doe      |
      | email     | john@test.com |
      | password  | 12345         |
    When I send a POST request to "/users"
    Then I should receive a 400 status code
    And the response should contain error message "password must be 6 characters or longer"

  Scenario: Register user with empty name
    Given I have invalid user data:
      | name      |               |
      | email     | john@test.com |
      | password  | password123   |
    When I send a POST request to "/users"
    Then I should receive a 400 status code
    And the response should contain error message "name mustn't be empty"

  Scenario: Register user with invalid JSON
    Given I send a POST request to "/users" with invalid JSON body
    Then I should receive a 400 status code
    And the response should contain error message "invalid json body"
