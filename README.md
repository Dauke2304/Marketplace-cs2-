# Marketplace CS2

## Description
Marketplace CS2 is a Go-based marketplace application designed for Counter-Strike 2 items. It provides a platform for users to buy, sell, and trade in-game items securely and efficiently.

## Features
- User authentication and account management
- Item listing with images and descriptions
- Admin panel for managing listings and users
- Search and filter functionality for easy navigation

## Installation
### Prerequisites
- Go (latest stable version)
- MongoDB

### Steps
1. Download .zip file
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Start the application:
   ```sh
   go run main.go
   ```

## Project structure
```
marketplace-cs2-/
├── cmd/          # Main application entry
├── config/       # Configuration files
├── controllers/  # HTTP controllers
├── database/     # Database connection
├── models/       # Data models
├── repositories/ # Database operations
├── services/     # Business logic
└── frontend/     # Frontend assets
```
### Implementation Notes:

1. Security: The admin panel checks for both valid session and the "admin" username
2. Data Management: Uses existing repositories for all DB operations
3. Navigation: Separate pages for users, skins, and transactions management
4. Styling: Consistent with your existing design but with admin-specific features

To use the admin panel:
1. Login with username "admin" and password "admin123"
2. Access `/admin` endpoint
3. Use navigation links to manage different entities

The admin panel uses your existing repositories (UserRepository, SkinRepository, TransactionRepository) for all operations, maintaining consistency with the rest of the application.

## Usage
1. Register or log in to your account.
2. Browse available items or list your own for sale.
3. Complete transactions securely.
4. Manage your account and listings via the dashboard.

## Authors
* Zhasulan Abdigaliyev
* Dauletkhan Bulanbayev
* Ernar Zhaksybaev